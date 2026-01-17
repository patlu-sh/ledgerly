package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimiter implements a simple token bucket rate limiter
type RateLimiter struct {
	mu       sync.Mutex
	clients  map[string]*client
	rate     int           // tokens per interval
	interval time.Duration // refill interval
	burst    int           // max tokens (bucket size)
}

type client struct {
	tokens    int
	lastCheck time.Time
}

// NewRateLimiter creates a rate limiter
// rate: requests allowed per interval
// interval: time period for rate
// burst: max requests allowed in burst
func NewRateLimiter(rate int, interval time.Duration, burst int) *RateLimiter {
	rl := &RateLimiter{
		clients:  make(map[string]*client),
		rate:     rate,
		interval: interval,
		burst:    burst,
	}
	
	// Cleanup old entries every minute
	go rl.cleanup()
	
	return rl
}

func (rl *RateLimiter) cleanup() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, c := range rl.clients {
			if time.Since(c.lastCheck) > 5*time.Minute {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	c, exists := rl.clients[ip]
	now := time.Now()

	if !exists {
		rl.clients[ip] = &client{
			tokens:    rl.burst - 1,
			lastCheck: now,
		}
		return true
	}

	// Calculate tokens to add based on elapsed time
	elapsed := now.Sub(c.lastCheck)
	tokensToAdd := int(elapsed / rl.interval) * rl.rate
	c.tokens += tokensToAdd
	if c.tokens > rl.burst {
		c.tokens = rl.burst
	}
	c.lastCheck = now

	if c.tokens > 0 {
		c.tokens--
		return true
	}

	return false
}

// Middleware returns a Gin middleware function
func (rl *RateLimiter) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		
		if !rl.allow(ip) {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "rate limit exceeded",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}

// DefaultRateLimiter returns a rate limiter with sensible defaults
// 100 requests per second, burst of 200
func DefaultRateLimiter() *RateLimiter {
	return NewRateLimiter(100, time.Second, 200)
}
