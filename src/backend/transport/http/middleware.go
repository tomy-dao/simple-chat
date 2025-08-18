package httpTransoprt

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
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

func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getToken(r)
		// Add token to request context
		fmt.Println("token", token)
		if token != "" {
			// Parse token to extract user_id and session_id
			claims, err := decodeJWT(token)
			if err != nil {
				http.Error(w, "Unauthorized: Token expired or invalid", http.StatusUnauthorized)
				return
			}
			
			ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
			ctx = context.WithValue(ctx, "session_id", claims["session_id"])
			ctx = context.WithValue(ctx, "token", token)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

func SetupMiddleware(r chi.Router) {
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.CleanPath)
	r.Use(middleware.GetHead)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Use(TokenMiddleware)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// JSON Content-Type middleware
	r.Use(JSONContentTypeMiddleware)
}

// JSONContentTypeMiddleware automatically sets Content-Type to application/json for all responses
func JSONContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set Content-Type to application/json for all responses
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware custom logging middleware
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		duration := time.Since(start)
		
		// Log request details
		// You can integrate with your logging system here
		_ = duration // Use duration to avoid unused variable warning
		// log.Printf("Request: %s %s - Duration: %v", r.Method, r.URL.Path, duration)
	})
}
