package limiter

import (
	"context"
	"smart-gateway/internal/clients"
	"time"
)

// 漏桶
type LeakyBucketLimiter struct {
	rdb      *clients.RedisDB
	rate     int // 每秒漏水数量
	capacity int // 最大队列容量
}

func NewLeakyBucketLimiter(rdb *clients.RedisDB, rate, capacity int) *LeakyBucketLimiter {
	return &LeakyBucketLimiter{rdb: rdb, rate: rate, capacity: capacity}
}

const leakyScript = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local rate = tonumber(ARGV[2])
local cap = tonumber(ARGV[3])

local last = tonumber(redis.call('hget', key, "ts") or now)
local water = tonumber(redis.call('hget', key, "water") or 0)

-- 流出水量
local leak = (now - last) / 1000 * rate
water = math.max(0, water - leak)

local pass = 0
if water + 1 <= cap then
    water = water + 1
    pass = 1
end

redis.call('hset', key, "water", water)
redis.call('hset', key, "ts", now)
redis.call('expire', key, 60)
return pass
`

func (l *LeakyBucketLimiter) Allow(ctx context.Context, key string) (bool, error) {
	nowMs := time.Now().UnixMilli()
	pass, err := l.rdb.Eval(ctx, leakyScript, []string{key}, nowMs, l.rate, l.capacity).Bool()
	return pass, err
}
