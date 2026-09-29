package limiter

import (
	"context"
	"smart-gateway/internal/clients"
	"time"
)

// 滑动窗口
type SlidingWindowLimiter struct {
	rdb    *clients.RedisDB
	rate   int
	window int
}

func NewSlidingWindowLimiter(rdb *clients.RedisDB, rate, window int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{rdb: rdb, rate: rate, window: window}
}

const slidingScript = `
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])
local expireTs = now - window * 1000
redis.call('zremrangebyscore', key, 0, expireTs)
local count = redis.call('zcard', key)
if count < limit then
    redis.call('zadd', key, now, now .. ':' .. math.random())
    redis.call('expire', key, window)
    return 1
end
return 0
`

func (l *SlidingWindowLimiter) Allow(ctx context.Context, key string) (bool, error) {
	nowMs := time.Now().UnixMilli()
	ok, err := l.rdb.Eval(ctx, slidingScript, []string{key}, nowMs, l.window, l.rate).Bool()
	return ok, err
}
