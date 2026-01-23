package middleware

import (
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"pet-everyone/cmd/web/application"

	"golang.org/x/time/rate"
)

// clientLimiter is a store of a client's rate limiter and the time they were last seen
type clientLimiter struct {
	ip       string
	limiter  *rate.Limiter
	lastSeen time.Time
}

// TODO: switch to map[string]*clientLimiter so we can use lastSeen for pruning
type rateLimiterMap struct {
	ipLimiterMap map[string]*rate.Limiter
	mu           sync.RWMutex
}

// prune iterates through active limiters and prunes TTL expired entries
func (r *rateLimiterMap) prune() {
	var interval time.Duration
	ticker := time.NewTicker(interval)
	for {
		select {
		case t := <-ticker.C:
			log.Println("Pruning rate limiter map at", t)
		}
	}
}

func RateLimit(app *application.Config, interval time.Duration, burst int) func(http.Handler) http.Handler {
	rlm := rateLimiterMap{
		ipLimiterMap: make(map[string]*rate.Limiter),
	}

	l := rate.Every(interval)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ip := getIP(r)

			rlm.mu.RLock()
			limiter, exists := rlm.ipLimiterMap[ip]

			if !exists {
				limiter = rate.NewLimiter(l, burst)
				rlm.ipLimiterMap[ip] = limiter
			}
			rlm.mu.RUnlock()

			if !limiter.Allow() {
				app.RespondWithError(w, http.StatusTooManyRequests, "Too many requests", nil)
				return
			}
			next.ServeHTTP(w, r)
		})

	}
}

func getIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Println(err)
		return ""
	}
	return host
}
