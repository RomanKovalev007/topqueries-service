package store

import "testing"

func TestRateLimiter_AllowUnderLimit(t *testing.T) {
	rl := NewRateLimiter()

	for range maxCountOfUserRequest {
		if !rl.Allow("user1") {
			t.Fatal("expected Allow=true under limit")
		}
	}
}

func TestRateLimiter_ExactlyAtLimit(t *testing.T) {
	rl := NewRateLimiter()

	for range maxCountOfUserRequest - 1 {
		rl.Allow("user1")
	}

	if !rl.Allow("user1") {
		t.Fatal("expected Allow=true exactly at limit")
	}
	if rl.Allow("user1") {
		t.Fatal("expected Allow=false one over limit")
	}
}

func TestRateLimiter_IndependentUsers(t *testing.T) {
	rl := NewRateLimiter()

	for range maxCountOfUserRequest + 1 {
		rl.Allow("user1")
	}

	if !rl.Allow("user2") {
		t.Fatal("expected user2 to be allowed regardless of user1 limit")
	}
}


