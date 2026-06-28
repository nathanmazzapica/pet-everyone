package middleware

import (
	"testing"
	"time"

	"golang.org/x/time/rate"
)

func testLimiterMap(t *testing.T) *rateLimiterMap {
	t.Helper()
	return &rateLimiterMap{
		ipLimiterMap:  make(map[string]*clientLimiter),
		limiterTTL:    1 * time.Millisecond,
		pruneInterval: 2 * time.Millisecond,
	}

}

func Test_RateLimitPrune(t *testing.T) {
	rlm := testLimiterMap(t)

	expired := &clientLimiter{
		ip:      "0.0.0.0",
		limiter: rate.NewLimiter(rate.Every(1*time.Minute), 15),
	}
	expired.touch(time.Now())

	safe := &clientLimiter{
		ip:      "0.0.0.1",
		limiter: rate.NewLimiter(rate.Every(1*time.Minute), 15),
	}
	safe.touch(time.Now().Add(5 * time.Minute))

	rlm.ipLimiterMap[expired.ip] = expired
	rlm.ipLimiterMap[safe.ip] = safe

	go rlm.prune()

	delay := time.NewTicker(5 * time.Millisecond)
	<-delay.C

	rlm.mu.RLock()
	defer rlm.mu.RUnlock()

	if _, ok := rlm.ipLimiterMap[expired.ip]; ok {
		t.Fatalf("Expected client1 to be pruned, was not")
	}

	if _, ok := rlm.ipLimiterMap[safe.ip]; !ok {
		t.Fatalf("Expected client2 to be present, was pruned")
	}
}
