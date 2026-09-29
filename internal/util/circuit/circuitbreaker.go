package circuit

import (
	"context"
	"strconv"
	"time"

	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
)

// State 熔断器状态。
type State string

const (
	StateClosed   State = "closed"    // 关闭，正常放行
	StateOpen     State = "open"      // 开路熔断，拒绝请求
	StateHalfOpen State = "half_open" // 半开探测
)

// CBreaker 基于 Redis 的分布式熔断器。
// 状态与计数均写入 Redis，保证网关多实例下熔断判定一致。
type CBreaker struct {
	rdb *clients.RedisDB

	name             string        // 熔断器标识（Redis key 后缀）
	failThreshold    int           // 连续失败阈值，达到即打开
	successThreshold int           // 半开状态成功阈值，达到即关闭
	resetTimeout     time.Duration // 打开后多久进入半开
	halfOpenMax      int           // 半开状态允许的并发探测请求数
	timeout          time.Duration // 下游调用超时
}

// NewCBreaker 根据熔断配置构造熔断器。
func NewCBreaker(rdb *clients.RedisDB, cb model.CircuitBreaker) *CBreaker {
	name := cb.Name
	if name == "" {
		name = strconv.FormatUint(cb.ID, 10)
	}
	return &CBreaker{
		rdb:              rdb,
		name:             name,
		failThreshold:    cb.FailureThreshold,
		successThreshold: cb.SuccessThreshold,
		resetTimeout:     time.Duration(cb.OpenDurationSeconds) * time.Second,
		halfOpenMax:      cb.HalfOpenMaxRequests,
		timeout:          time.Duration(cb.TimeoutMs) * time.Millisecond,
	}
}

func (cb *CBreaker) stateKey() string    { return "circuit:" + cb.name + ":state" }
func (cb *CBreaker) openTimeKey() string { return "circuit:" + cb.name + ":open_time" }
func (cb *CBreaker) probeKey() string    { return "circuit:" + cb.name + ":probe" }
func (cb *CBreaker) failKey() string     { return "circuit:" + cb.name + ":fail" }
func (cb *CBreaker) successKey() string  { return "circuit:" + cb.name + ":success" }

// allowScript 判定是否放行并推进状态机。
// KEYS[1]=state KEYS[2]=open_time KEYS[3]=probe KEYS[4]=fail KEYS[5]=success
// ARGV[1]=now_ms ARGV[2]=reset_timeout_ms ARGV[3]=half_open_max
const allowScript = `
local state = redis.call('get', KEYS[1]) or 'closed'
local now = tonumber(ARGV[1])
local resetMs = tonumber(ARGV[2])
local halfOpenMax = tonumber(ARGV[3])

if state == 'closed' then
    return 1
end

if state == 'open' then
    local openTime = tonumber(redis.call('get', KEYS[2]) or 0)
    if now - openTime >= resetMs then
        redis.call('set', KEYS[1], 'half_open')
        redis.call('set', KEYS[5], 0)
        redis.call('incr', KEYS[3])
        redis.call('expire', KEYS[3], 60)
        return 1
    end
    return 0
end

if state == 'half_open' then
    local probe = tonumber(redis.call('get', KEYS[3]) or 0)
    if probe < halfOpenMax then
        redis.call('incr', KEYS[3])
        redis.call('expire', KEYS[3], 60)
        return 1
    end
    return 0
end

return 1
`

// successScript 请求成功回调。
// KEYS[1]=state KEYS[2]=open_time KEYS[3]=probe KEYS[4]=fail KEYS[5]=success
// ARGV[1]=success_threshold
const successScript = `
local state = redis.call('get', KEYS[1]) or 'closed'

if state == 'closed' then
    redis.call('del', KEYS[4])
    return
end

if state == 'half_open' then
    local probe = tonumber(redis.call('get', KEYS[3]) or 0)
    if probe > 0 then
        redis.call('decr', KEYS[3])
    end
    local success = redis.call('incr', KEYS[5])
    redis.call('expire', KEYS[5], 60)
    if success >= tonumber(ARGV[1]) then
        redis.call('set', KEYS[1], 'closed')
        redis.call('del', KEYS[3])
        redis.call('del', KEYS[4])
        redis.call('del', KEYS[5])
    end
end
`

// failScript 请求失败回调（超时、5xx 均视为失败）。
// KEYS[1]=state KEYS[2]=open_time KEYS[3]=probe KEYS[4]=fail KEYS[5]=success
// ARGV[1]=now_ms ARGV[2]=fail_threshold
const failScript = `
local state = redis.call('get', KEYS[1]) or 'closed'
local now = tonumber(ARGV[1])

if state == 'closed' then
    local fail = redis.call('incr', KEYS[4])
    redis.call('expire', KEYS[4], 60)
    if fail >= tonumber(ARGV[2]) then
        redis.call('set', KEYS[1], 'open')
        redis.call('set', KEYS[2], now)
        redis.call('del', KEYS[3])
    end
    return
end

if state == 'half_open' then
    redis.call('set', KEYS[1], 'open')
    redis.call('set', KEYS[2], now)
    redis.call('del', KEYS[3])
    redis.call('del', KEYS[5])
end
`

// Allow 判断本次请求是否放行；返回 false 直接熔断拒绝。
func (cb *CBreaker) Allow(ctx context.Context) bool {
	pass, err := cb.rdb.Eval(ctx, allowScript,
		[]string{cb.stateKey(), cb.openTimeKey(), cb.probeKey(), cb.failKey(), cb.successKey()},
		time.Now().UnixMilli(), cb.resetTimeout.Milliseconds(), cb.halfOpenMax,
	).Int()
	if err != nil {
		// Redis 异常时按放行处理，避免熔断器自身故障拖垮整条链路。
		return true
	}
	return pass == 1
}

// OnSuccess 请求下游成功之后回调。
func (cb *CBreaker) OnSuccess(ctx context.Context) {
	_ = cb.rdb.Eval(ctx, successScript,
		[]string{cb.stateKey(), cb.openTimeKey(), cb.probeKey(), cb.failKey(), cb.successKey()},
		cb.successThreshold,
	).Err()
}

// OnFail 请求下游失败之后回调。
func (cb *CBreaker) OnFail(ctx context.Context) {
	_ = cb.rdb.Eval(ctx, failScript,
		[]string{cb.stateKey(), cb.openTimeKey(), cb.probeKey(), cb.failKey(), cb.successKey()},
		time.Now().UnixMilli(), cb.failThreshold,
	).Err()
}

// State 获取当前状态，用于日志监控。
func (cb *CBreaker) State(ctx context.Context) State {
	s, err := cb.rdb.Get(ctx, cb.stateKey()).Result()
	if err != nil || s == "" {
		return StateClosed
	}
	return State(s)
}

// Timeout 下游调用超时。
func (cb *CBreaker) Timeout() time.Duration { return cb.timeout }

// Reset 手动复位熔断器到关闭状态。
func (cb *CBreaker) Reset(ctx context.Context) {
	_ = cb.rdb.Del(ctx, cb.stateKey(), cb.openTimeKey(), cb.probeKey(), cb.failKey(), cb.successKey()).Err()
}
