package httpTransoprt

import (
	"encoding/json"
	"local/endpoint"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func handleRouter(r chi.Router, endpoints *endpoint.Endpoints) chi.Router {
	h := &handler{endpoints: endpoints}

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth endpoints
		r.Route("/auth", func(r chi.Router) {
			r.Get("/me", h.GetMe())
			r.Post("/register", h.Register())
			r.Post("/login", h.Login())
			r.Post("/logout", h.Logout())
		})

		// User endpoints
		r.Route("/users", func(r chi.Router) {
			r.Get("/", h.GetUsers())
		})

		// Conversation endpoints
		r.Route("/conversations", func(r chi.Router) {
			r.Post("/", h.CreateConversation())
			r.Get("/", h.GetConversations())
			r.Get("/user/{userID}", h.GetConversationByUserID())
			r.Post("/{conversationID}/messages", h.CreateMessage())
			r.Get("/{conversationID}/messages", h.GetMessagesByConversationID())
		})
	})

	// Health check
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"status":    "OK",
			"timestamp": time.Now().UTC().Format(time.RFC3339),
			"service":   "local-service",
		}
		json.NewEncoder(w).Encode(response)
	})

	// Root endpoint
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		response := map[string]interface{}{
			"message": "Local Service API",
			"version": "1.0.0",
			"endpoints": map[string]string{
				"health":     "/health",
				"register":   "/api/v1/auth/register",
				"login":      "/api/v1/auth/login",
				"logout":     "/api/v1/auth/logout",
				"conversations": "/api/v1/conversations",
			},
		}
		json.NewEncoder(w).Encode(response)
	})

	return r
}
