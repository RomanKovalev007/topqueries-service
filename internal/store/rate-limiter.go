package store

import (
	"context"
	"time"
)

const (
	maxCountOfUserRequest = 100
	rateLimiterLengthOftime = 1*time.Minute
	rateLimiterCountsOfBuckets = 12
)

type RateLimiter struct {
	window *CounterWindow
}

func NewRateLimiter() *RateLimiter {
	window := NewCounterWindow(rateLimiterLengthOftime, rateLimiterCountsOfBuckets)
	return &RateLimiter{
		window: window,
	}
}

func (rl *RateLimiter) Allow(userID string) bool {
	rl.window.Add(userID)
	return rl.window.GetCount(userID) <= maxCountOfUserRequest
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