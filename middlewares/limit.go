package middlewares

import (
	"auth-service/utils"
	"net/http"
	"sync"
	"time"
)

// Define a struct to hold rate limiting data
type rateLimiteType struct {
	rate       int           // requests per interval
	interval   time.Duration // interval duration
	tokens     int
	lastRefill time.Time
	mu         sync.Mutex
}

func NewRateLimiter(rate int, interval time.Duration) *rateLimiteType {
	return &rateLimiteType{
		rate:       rate,
		interval:   interval,
		tokens:     rate,
		lastRefill: time.Now(),
	}
}

func (rl *rateLimiteType) allow() bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(rl.lastRefill)

	// Refill tokens based on elapsed time
	if elapsed >= rl.interval {
		rl.tokens = rl.rate
		rl.lastRefill = now
	} else {
		refill := int(elapsed / rl.interval)
		rl.tokens = utils.Min(rl.rate, rl.tokens+refill)
	}

	// Check if there are enough tokens available
	if rl.tokens > 0 {
		rl.tokens--
		return true
	}

	return false
}

// LimitRate creates a middleware that applies rate limiting
func LimitRate(next http.HandlerFunc, rl *rateLimiteType) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if rl.allow() {
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
		}
	}
}
