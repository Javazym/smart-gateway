package limiter

import (
	"context"
	"smart-gateway/internal/clients"
	"time"
)

// 令牌桶
type TokenBucketLimiter struct {
	rdb   *clients.RedisDB
	rate  int // 每秒补充令牌数
	burst int // 桶最大容量
}

func NewTokenBucketLimiter(rdb *clients.RedisDB, rate, burst int) *TokenBucketLimiter {
	return &TokenBucketLimiter{rdb: rdb, rate: rate, burst: burst}
}

const tokenBucketScript = `
local keyTokens = KEYS[1]
local keyTs = KEYS[2]
local now = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local burst = tonumber(ARGV[3])

local tokens = tonumber(redis.call('get', keyTokens) or burst)
local lastTs = tonumber(redis.call('get', keyTs) or now)

-- 计算新增令牌
local delta = math.max(0, now - lastTs)
local add = delta / 1000 * rate
tokens = math.min(burst, tokens + add)

local pass = 0
if tokens >= 1 then
    tokens = tokens - 1
    pass = 1
end

redis.call('set', keyTokens, tokens)
redis.call('set', keyTs, now)
redis.call('expire', keyTokens, 60)
redis.call('expire', keyTs, 60)
return pass
`

func (l *TokenBucketLimiter) Allow(ctx context.Context, key string) (bool, error) {
	nowMs := time.Now().UnixMilli()
	tokenKey := key + ":token"
	tsKey := key + ":ts"
	pass, err := l.rdb.Eval(ctx, tokenBucketScript, []string{tokenKey, tsKey}, nowMs, l.rate, l.burst).Bool()
	return pass, err
}
