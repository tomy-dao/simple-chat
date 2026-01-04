package httpTransport

import (
	"local/config"
	"local/endpoint"
	"local/model"

	"github.com/gin-gonic/gin"
)

func MakeHttpTransport(initParams *model.InitParams, endpoints *endpoint.Endpoints) *gin.Engine {
	r := gin.Default()

	// Initialize rate limiter if enabled
	if config.Config.RateLimitEnabled {
		InitRateLimiter(config.Config.RateLimitRequestsPerMin, config.Config.RateLimitBurst)
	}

	SetupMiddleware(r)

	handleRouter(r, endpoints)

	return r
}
