package httpTransoprt

import (
	"local/endpoint"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func handleRouter(r *gin.Engine, endpoints *endpoint.Endpoints) {
	h := &handler{endpoints: endpoints}

	// Swagger documentation
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API routes
	v1 := r.Group("/api/v1")
	{
		// Auth endpoints
		v1.POST("/register", h.Register())
		v1.POST("/login", h.Login())

		// Protected routes
		protected := v1.Group("")
		protected.Use(ProtectedMiddleware(endpoints))
		{
			protected.POST("/logout", h.Logout())
			protected.GET("/me", h.GetMe())

			// Users endpoints
			users := protected.Group("/users")
			{
				users.GET("/", h.GetUsers())
			}

			// Conversation endpoints
			conversations := protected.Group("/conversations")
			{
				conversations.POST("/", h.CreateConversation())
				conversations.GET("/", h.GetConversations())
				conversations.GET("/user/:userID", h.GetConversationByUserID())
				conversations.POST("/:conversationID/messages", h.CreateMessage())
				conversations.GET("/:conversationID/messages", h.GetMessagesByConversationID())
			}
		}
	}

	// Health check
	r.GET("/health", func(c *gin.Context) {
		OK(c, gin.H{
			"status":    "OK",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"service":   "local-service",
		})
	})

	// Root endpoint
	r.GET("/", func(c *gin.Context) {
		OK(c, gin.H{
			"message": "Local Service API",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"health":        "/health",
				"swagger":       "/swagger/index.html",
				"register":      "/api/v1/register",
				"login":         "/api/v1/login",
				"logout":        "/api/v1/logout",
				"conversations": "/api/v1/conversations",
			},
		}, "API is running")
	})
}
