package store

import (
	"testing"
	"time"
)

const testMaxRequests = int64(100)

func newTestRateLimiter() *RateLimiter {
	return NewRateLimiter(time.Minute, 12, testMaxRequests)
}

func TestRateLimiter_AllowUnderLimit(t *testing.T) {
	rl := newTestRateLimiter()

	for range testMaxRequests {
		if !rl.Allow("user1") {
			t.Fatal("expected Allow=true under limit")
		}
	}
}

func TestRateLimiter_BlockOverLimit(t *testing.T) {
	rl := newTestRateLimiter()

	for range testMaxRequests {
		rl.Allow("user1")
	}

	if rl.Allow("user1") {
		t.Fatal("expected Allow=false after exceeding limit")
	}
}

func TestRateLimiter_IndependentUsers(t *testing.T) {
	rl := newTestRateLimiter()

	for range testMaxRequests + 1 {
		rl.Allow("user1")
	}

	if !rl.Allow("user2") {
		t.Fatal("expected user2 to be allowed regardless of user1 limit")
	}
}

func TestRateLimiter_ExactlyAtLimit(t *testing.T) {
	rl := newTestRateLimiter()

	for range testMaxRequests - 1 {
		rl.Allow("user1")
	}

	if !rl.Allow("user1") {
		t.Fatal("expected Allow=true exactly at limit")
	}
	if rl.Allow("user1") {
		t.Fatal("expected Allow=false one over limit")
	}
}
