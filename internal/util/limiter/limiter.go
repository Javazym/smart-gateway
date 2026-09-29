package limiter

import (
	"context"
	"smart-gateway/internal/clients"
	"smart-gateway/internal/model"
)

type Limiter interface {
	Allow(ctx context.Context, key string) (bool, error)
}

func NewLimiter(rdb *clients.RedisDB, policy *model.RateLimitPolicy) Limiter {
	switch policy.Algorithm {
	case "fixed_window":
		return NewFixedWindowLimiter(rdb, policy.Rate, policy.WindowSeconds)
	case "sliding_window":
		return NewSlidingWindowLimiter(rdb, policy.Rate, policy.WindowSeconds)
	case "token_bucket":
		burst := policy.Burst
		if burst <= 0 {
			burst = policy.Rate
		}
		return NewTokenBucketLimiter(rdb, policy.Rate, burst)
	case "leaky_bucket":
		capacity := policy.Burst
		if capacity <= 0 {
			capacity = policy.Rate
		}
		return NewLeakyBucketLimiter(rdb, policy.Rate, capacity)
	default:
		return nil
	}
}
