package httpTransport

import (
	"net"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// visitor tracks a rate limiter and last access time for cleanup
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiterManager manages rate limiters for different keys (user_id or IP)
type RateLimiterManager struct {
	visitors map[string]*visitor
	mu       sync.RWMutex
	rate     rate.Limit // requests per second
	burst    int        // maximum burst capacity
}

// NewRateLimiterManager creates a new rate limiter manager
func NewRateLimiterManager(requestsPerMin int, burst int) *RateLimiterManager {
	rlm := &RateLimiterManager{
		visitors: make(map[string]*visitor),
		rate:     rate.Limit(float64(requestsPerMin) / 60.0), // convert per minute to per second
		burst:    burst,
	}

	// Start cleanup goroutine to prevent memory leaks
	go rlm.cleanupVisitors()

	return rlm
}

// GetLimiter retrieves or creates a rate limiter for the given key
func (rlm *RateLimiterManager) GetLimiter(key string) *rate.Limiter {
	rlm.mu.Lock()
	defer rlm.mu.Unlock()

	v, exists := rlm.visitors[key]
	if !exists {
		limiter := rate.NewLimiter(rlm.rate, rlm.burst)
		rlm.visitors[key] = &visitor{
			limiter:  limiter,
			lastSeen: time.Now(),
		}
		return limiter
	}

	// Update last seen time
	v.lastSeen = time.Now()
	return v.limiter
}

// cleanupVisitors removes inactive visitors to prevent memory leaks
func (rlm *RateLimiterManager) cleanupVisitors() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		rlm.mu.Lock()
		for key, v := range rlm.visitors {
			// Remove visitors not seen in last 10 minutes
			if time.Since(v.lastSeen) > 10*time.Minute {
				delete(rlm.visitors, key)
			}
		}
		rlm.mu.Unlock()
	}
}

// getClientIP extracts the client IP address from the request
// It checks X-Forwarded-For, X-Real-IP headers, and falls back to RemoteAddr
func getClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header first (for proxies/load balancers)
	xff := c.GetHeader("X-Forwarded-For")
	if xff != "" {
		// X-Forwarded-For can contain multiple IPs, take the first one
		if ip, _, err := net.SplitHostPort(xff); err == nil {
			return ip
		}
		return xff
	}

	// Check X-Real-IP header
	xri := c.GetHeader("X-Real-IP")
	if xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	if ip, _, err := net.SplitHostPort(c.Request.RemoteAddr); err == nil {
		return ip
	}

	return c.Request.RemoteAddr
}
