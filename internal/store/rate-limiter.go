package store

import (
	"context"
	"time"
)

type RateLimiter struct {
	window      *CounterWindow
	maxRequests int64
}

func NewRateLimiter(window time.Duration, bucketCount int, maxRequests int64) *RateLimiter {
	return &RateLimiter{
		window:      NewCounterWindow(window, bucketCount),
		maxRequests: maxRequests,
	}
}

func (rl *RateLimiter) Allow(userID string) bool {
	rl.window.Add(userID)
	return rl.window.GetCount(userID) <= rl.maxRequests
}

func (rl *RateLimiter) Ticker(ctx context.Context) {
	ticker := time.NewTicker(rl.window.GetBucketInterval())
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			rl.window.Rotate()
		}
	}
}
