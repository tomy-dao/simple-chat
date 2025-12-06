package httpTransoprt

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"local/endpoint"
	"local/model"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func decodeJWT(tokenString string) (map[string]interface{}, error) {
	// Split the token into parts
	parts := strings.Split(tokenString, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid JWT format")
	}

	// Decode the payload (second part)
	payload := parts[1]
	
	// Add padding if necessary
	if len(payload)%4 != 0 {
		payload += strings.Repeat("=", 4-len(payload)%4)
	}

	// Base64 decode
	decoded, err := base64.URLEncoding.DecodeString(payload)
	if err != nil {
		return nil, err
	}

	// Parse JSON
	var claims map[string]interface{}
	err = json.Unmarshal(decoded, &claims)
	if err != nil {
		return nil, err
	}

	return claims, nil
}

func getToken(c *gin.Context) string {
	token := c.GetHeader("Authorization")
	if token == "" {
		return ""
	}
	if len(token) > 7 && token[:7] == "Bearer " {
		return token[7:]
	}
	return token
}

func TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := getToken(c)
		// Add token to request context
		if token != "" {
			// Parse token to extract user_id and session_id
			claims, err := decodeJWT(token)
			if err != nil {
				// Don't abort here, just skip adding context
				// Let ProtectedMiddleware handle authentication
			} else {
				var userID uint
				if uid, ok := claims["user_id"].(float64); ok {
					userID = uint(uid)
				}
				sessionID := ""
				if sid, ok := claims["session_id"].(string); ok {
					sessionID = sid
				}
				reqCtx := model.NewRequestContext(c.Request.Context()).WithClaims(token, userID, sessionID)
				c.Request = c.Request.WithContext(reqCtx.Context())
			}
		}

		c.Next()
	}
}

func SetupMiddleware(r *gin.Engine) {
	// OpenTelemetry tracing middleware (must be first)
	r.Use(otelgin.Middleware("simple-chat-api"))

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"}
	config.ExposeHeaders = []string{"Link"}
	config.AllowCredentials = true
	config.MaxAge = 12 * time.Hour
	r.Use(cors.New(config))

	// Token middleware
	r.Use(TokenMiddleware())

	// JSON Content-Type middleware
	r.Use(JSONContentTypeMiddleware())
}

// JSONContentTypeMiddleware automatically sets Content-Type to application/json for all responses
// Excludes Swagger routes to allow HTML content type
func JSONContentTypeMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip JSON content type for Swagger routes
		if strings.HasPrefix(c.Request.URL.Path, "/swagger") {
			c.Next()
			return
		}
		c.Header("Content-Type", "application/json")
		c.Next()
	}
}

func ProtectedMiddleware(endpoints *endpoint.Endpoints) gin.HandlerFunc {
	return func(c *gin.Context) {
		reqCtx := model.NewRequestContext(c.Request.Context())
		response := endpoints.Auth.Authenticate(reqCtx)
		if !response.OK() {
			Unauthorized(c, response.ErrorString())
			c.Abort()
			return
		}
		c.Next()
	}
}
