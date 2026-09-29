package limiter

import (
	"context"
	"smart-gateway/internal/clients"
)

// 固定窗口
type FixedWindowLimiter struct {
	rdb    *clients.RedisDB
	rate   int
	window int
}

func NewFixedWindowLimiter(rdb *clients.RedisDB, rate, window int) *FixedWindowLimiter {
	return &FixedWindowLimiter{rdb: rdb, rate: rate, window: window}
}

const fixedWindowScript = `
local key = KEYS[1]
local limit = tonumber(ARGV[1])
local expire = tonumber(ARGV[2])
local current = redis.call('incr', key)
if current == 1 then
redis.call('expire', key, expire)
end
if current <= limit then
    return 1
end
return 0
`

func (l *FixedWindowLimiter) Allow(ctx context.Context, key string) (bool, error) {
	res, err := l.rdb.Eval(ctx, fixedWindowScript, []string{key}, l.rate, l.window).Bool()
	if err != nil {
		return false, err
	}
	return res, nil
}
