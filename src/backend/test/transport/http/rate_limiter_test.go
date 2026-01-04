package http_test

import (
	"testing"
	"time"

	httpTransport "local/transport/http"

	"github.com/stretchr/testify/assert"
)

func TestRateLimiterManager_GetLimiter(t *testing.T) {
	rlm := httpTransport.NewRateLimiterManager(60, 10)

	limiter1 := rlm.GetLimiter("user:1")
	limiter2 := rlm.GetLimiter("user:1")

	// Same key should return same limiter
	assert.Equal(t, limiter1, limiter2, "Same key should return same limiter instance")
}

func TestRateLimiterManager_PerUserIsolation(t *testing.T) {
	rlm := httpTransport.NewRateLimiterManager(60, 2) // 2 requests burst

	limiter1 := rlm.GetLimiter("user:1")
	limiter2 := rlm.GetLimiter("user:2")

	// Exhaust user1's quota
	assert.True(t, limiter1.Allow(), "User1 request 1 should be allowed")
	assert.True(t, limiter1.Allow(), "User1 request 2 should be allowed")
	assert.False(t, limiter1.Allow(), "User1 request 3 should be denied")

	// User2 should still have quota
	assert.True(t, limiter2.Allow(), "User2 request 1 should be allowed")
	assert.True(t, limiter2.Allow(), "User2 request 2 should be allowed")
}

func TestRateLimiterManager_BurstCapacity(t *testing.T) {
	rlm := httpTransport.NewRateLimiterManager(60, 5)
	limiter := rlm.GetLimiter("user:1")

	// Should allow burst of 5 requests
	for i := 0; i < 5; i++ {
		assert.True(t, limiter.Allow(), "Request %d should be allowed", i+1)
	}

	// 6th request should be denied
	assert.False(t, limiter.Allow(), "6th request should be denied")
}

func TestRateLimiterManager_TokenRefill(t *testing.T) {
	rlm := httpTransport.NewRateLimiterManager(60, 2) // 1 request per second
	limiter := rlm.GetLimiter("user:1")

	// Exhaust quota
	assert.True(t, limiter.Allow(), "Request 1 should be allowed")
	assert.True(t, limiter.Allow(), "Request 2 should be allowed")
	assert.False(t, limiter.Allow(), "Request 3 should be denied immediately")

	// Wait for token refill (1 second = 1 token for 60 req/min)
	time.Sleep(1100 * time.Millisecond)

	// Should allow 1 more request
	assert.True(t, limiter.Allow(), "Request after refill should be allowed")
}

func TestRateLimiterManager_MultipleVisitors(t *testing.T) {
	rlm := httpTransport.NewRateLimiterManager(60, 3)

	// Create multiple visitors and verify they each get their own limiter
	limiter1 := rlm.GetLimiter("user:1")
	limiter2 := rlm.GetLimiter("user:2")
	limiter3 := rlm.GetLimiter("ip:192.168.1.1")

	// Verify all limiters are distinct by exhausting one and checking others still work
	assert.True(t, limiter1.Allow())
	assert.True(t, limiter1.Allow())
	assert.True(t, limiter1.Allow())
	assert.False(t, limiter1.Allow(), "User1 should be rate limited")

	// Other limiters should still work
	assert.True(t, limiter2.Allow(), "User2 should not be affected")
	assert.True(t, limiter3.Allow(), "IP should not be affected")
}

func TestGetClientIP_XForwardedFor(t *testing.T) {
	// This test would require creating a gin.Context with headers
	// For now, we'll skip implementation details
	// In a real test, you'd create a mock gin.Context
	t.Skip("Requires gin.Context setup - covered in integration tests")
}
