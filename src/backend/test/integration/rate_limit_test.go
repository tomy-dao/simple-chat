package integration

import (
	"local/config"
	httpTransport "local/transport/http"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRateLimit_AuthenticatedUser(t *testing.T) {
	setup, err := SetupTestEnvironment()
	assert.NoError(t, err)
	defer setup.CleanupTestEnvironment()

	// Enable rate limiting with low limit for testing
	config.Config.RateLimitEnabled = true
	config.Config.RateLimitRequestsPerMin = 60
	config.Config.RateLimitBurst = 3
	httpTransport.InitRateLimiter(config.Config.RateLimitRequestsPerMin, config.Config.RateLimitBurst)

	token := registerAndLogin(t, setup, "testuser", "password123")

	// First 3 requests should succeed
	for i := 0; i < 3; i++ {
		recorder := makeRequest(setup, http.MethodGet, "/api/v1/me", token, nil)
		assert.Equal(t, http.StatusOK, recorder.Code, "Request %d should succeed", i+1)
		assert.NotEmpty(t, recorder.Header().Get("X-RateLimit-Limit"), "Should have X-RateLimit-Limit header")
		assert.NotEmpty(t, recorder.Header().Get("X-RateLimit-Remaining"), "Should have X-RateLimit-Remaining header")
	}

	// 4th request should be rate limited
	recorder := makeRequest(setup, http.MethodGet, "/api/v1/me", token, nil)
	assert.Equal(t, http.StatusTooManyRequests, recorder.Code, "4th request should be rate limited")
	assert.Equal(t, "0", recorder.Header().Get("X-RateLimit-Remaining"), "Remaining should be 0")
	assert.NotEmpty(t, recorder.Header().Get("Retry-After"), "Should have Retry-After header")

	response := parseResponse[interface{}](t, recorder)
	assert.False(t, response.OK(), "Response should not be OK")
}

func TestRateLimit_UnauthenticatedUser(t *testing.T) {
	setup, err := SetupTestEnvironment()
	assert.NoError(t, err)
	defer setup.CleanupTestEnvironment()

	config.Config.RateLimitEnabled = true
	config.Config.RateLimitBurst = 2
	httpTransport.InitRateLimiter(2, 2)

	// Unauthenticated requests to /login should be rate limited by IP
	for i := 0; i < 2; i++ {
		body := map[string]string{"username": "test", "password": "test"}
		recorder := makeRequest(setup, http.MethodPost, "/api/v1/login", "", body)
		// Don't care if login fails, just testing rate limit
		assert.NotEqual(t, http.StatusTooManyRequests, recorder.Code, "Request %d should not be rate limited", i+1)
	}

	// 3rd request should be rate limited
	body := map[string]string{"username": "test", "password": "test"}
	recorder := makeRequest(setup, http.MethodPost, "/api/v1/login", "", body)
	assert.Equal(t, http.StatusTooManyRequests, recorder.Code, "3rd request should be rate limited")
}

func TestRateLimit_PerUserIsolation(t *testing.T) {
	setup, err := SetupTestEnvironment()
	assert.NoError(t, err)
	defer setup.CleanupTestEnvironment()

	// Register users first with high limits
	config.Config.RateLimitEnabled = false
	user1Token := registerAndLogin(t, setup, "user1", "password123")
	user2Token := registerAndLogin(t, setup, "user2", "password123")

	// Now enable rate limiting with low limits for testing
	config.Config.RateLimitEnabled = true
	config.Config.RateLimitBurst = 2
	httpTransport.InitRateLimiter(60, 2)

	// Exhaust user1's quota
	for i := 0; i < 2; i++ {
		recorder := makeRequest(setup, http.MethodGet, "/api/v1/me", user1Token, nil)
		assert.Equal(t, http.StatusOK, recorder.Code, "User1 request %d should succeed", i+1)
	}
	recorder := makeRequest(setup, http.MethodGet, "/api/v1/me", user1Token, nil)
	assert.Equal(t, http.StatusTooManyRequests, recorder.Code, "User1 should be rate limited")

	// User2 should still be able to make requests
	recorder = makeRequest(setup, http.MethodGet, "/api/v1/me", user2Token, nil)
	assert.Equal(t, http.StatusOK, recorder.Code, "User2 should not be affected by user1's rate limit")
}

func TestRateLimit_ExcludedEndpoints(t *testing.T) {
	setup, err := SetupTestEnvironment()
	assert.NoError(t, err)
	defer setup.CleanupTestEnvironment()

	config.Config.RateLimitEnabled = true
	config.Config.RateLimitBurst = 0 // Set to 0 to ensure rate limiting would occur
	httpTransport.InitRateLimiter(0, 0)

	// Health check should never be rate limited
	for i := 0; i < 10; i++ {
		recorder := makeRequest(setup, http.MethodGet, "/health", "", nil)
		assert.Equal(t, http.StatusOK, recorder.Code, "Health check request %d should succeed", i+1)
		assert.Empty(t, recorder.Header().Get("X-RateLimit-Limit"), "Health endpoint should not have rate limit headers")
	}
}

func TestRateLimit_Disabled(t *testing.T) {
	setup, err := SetupTestEnvironment()
	assert.NoError(t, err)
	defer setup.CleanupTestEnvironment()

	// Disable rate limiting
	config.Config.RateLimitEnabled = false

	token := registerAndLogin(t, setup, "testuser", "password123")

	// Should be able to make many requests without rate limiting
	for i := 0; i < 20; i++ {
		recorder := makeRequest(setup, http.MethodGet, "/api/v1/me", token, nil)
		assert.Equal(t, http.StatusOK, recorder.Code, "Request %d should succeed when rate limiting is disabled", i+1)
	}
}

func TestRateLimit_ResponseHeaders(t *testing.T) {
	setup, err := SetupTestEnvironment()
	assert.NoError(t, err)
	defer setup.CleanupTestEnvironment()

	config.Config.RateLimitEnabled = true
	config.Config.RateLimitRequestsPerMin = 60
	config.Config.RateLimitBurst = 5
	httpTransport.InitRateLimiter(60, 5)

	token := registerAndLogin(t, setup, "testuser", "password123")

	recorder := makeRequest(setup, http.MethodGet, "/api/v1/me", token, nil)
	assert.Equal(t, http.StatusOK, recorder.Code)

	// Verify rate limit headers are present
	assert.NotEmpty(t, recorder.Header().Get("X-RateLimit-Limit"), "Should have X-RateLimit-Limit header")
	assert.NotEmpty(t, recorder.Header().Get("X-RateLimit-Remaining"), "Should have X-RateLimit-Remaining header")
	assert.NotEmpty(t, recorder.Header().Get("X-RateLimit-Reset"), "Should have X-RateLimit-Reset header")
}
