package middleware

import (
	"net"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"golang.org/x/time/rate"
	"pet-everyone/cmd/web/application"
)

const (
	DefaultLimiterTTL    = 30 * time.Minute
	DefaultPruneInterval = 5 * time.Minute
)

// clientLimiter is a store of a client's rate limiter and the time they were last seen.
// lastSeen is stored as unix nanos in an atomic so the request hot path can update it
// without taking rateLimiterMap's write lock.
type clientLimiter struct {
	ip       string
	limiter  *rate.Limiter
	lastSeen atomic.Int64
}

// touch records the last time this client was seen.
func (c *clientLimiter) touch(t time.Time) {
	c.lastSeen.Store(t.UnixNano())
}

type rateLimiterMap struct {
	ipLimiterMap  map[string]*clientLimiter
	mu            sync.RWMutex
	limiterTTL    time.Duration
	pruneInterval time.Duration
}

// prune iterates through active limiters and prunes TTL expired entries
func (r *rateLimiterMap) prune() {
	ticker := time.NewTicker(r.pruneInterval)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()

		r.mu.Lock()
		for ip, lim := range r.ipLimiterMap {
			lastSeen := time.Unix(0, lim.lastSeen.Load())
			if now.After(lastSeen.Add(r.limiterTTL)) {
				delete(r.ipLimiterMap, ip)
			}
		}
		r.mu.Unlock()
	}
}

// getClient returns the limiter for ip, creating it if absent, and records the
// client as just seen. Creation is double-checked under the write lock so that
// concurrent first-requests from the same IP share a single limiter.
func (r *rateLimiterMap) getClient(ip string, l rate.Limit, burst int) *clientLimiter {
	now := time.Now()

	r.mu.RLock()
	client, exists := r.ipLimiterMap[ip]
	r.mu.RUnlock()

	if exists {
		client.touch(now)
		return client
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	// Another goroutine may have created the client between releasing the read
	// lock and acquiring the write lock.
	if client, exists = r.ipLimiterMap[ip]; exists {
		client.touch(now)
		return client
	}

	client = &clientLimiter{
		ip:      ip,
		limiter: rate.NewLimiter(l, burst),
	}
	client.touch(now)
	r.ipLimiterMap[ip] = client

	return client
}

func RateLimit(app *application.Config, interval time.Duration, burst int) func(http.Handler) http.Handler {
	rlm := rateLimiterMap{
		ipLimiterMap:  make(map[string]*clientLimiter),
		limiterTTL:    DefaultLimiterTTL,
		pruneInterval: DefaultPruneInterval,
	}

	l := rate.Every(interval)

	go rlm.prune()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip, err := getIP(r)

			if err != nil {
				app.RespondWithError(w, http.StatusInternalServerError, "An unknown error occurred", nil)
				return
			}

			client := rlm.getClient(ip, l, burst)

			if !client.limiter.Allow() {
				app.RespondWithError(w, http.StatusTooManyRequests, "Too many requests", nil)
				return
			}
			next.ServeHTTP(w, r)
		})

	}
}

func getIP(r *http.Request) (string, error) {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return "", err
	}
	return host, nil
}
