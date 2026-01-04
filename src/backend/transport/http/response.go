package httpTransport

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Standard API Response structure
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
	Message string      `json:"message,omitempty"`
}

// Standard API Error structure
type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Error codes
const (
	ErrCodeValidation        = "VALIDATION_ERROR"
	ErrCodeUnauthorized      = "UNAUTHORIZED"
	ErrCodeForbidden         = "FORBIDDEN"
	ErrCodeNotFound          = "NOT_FOUND"
	ErrCodeInternal          = "INTERNAL_ERROR"
	ErrCodeBadRequest        = "BAD_REQUEST"
	ErrCodeConflict          = "CONFLICT"
	ErrCodeRateLimitExceeded = "RATE_LIMIT_EXCEEDED"
)

// SuccessResponse returns a successful response
func SuccessResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	response := APIResponse{
		Success: true,
		Data:    data,
		Message: message,
	}
	c.JSON(statusCode, response)
}

// ErrorResponse returns an error response
func ErrorResponse(c *gin.Context, statusCode int, code, message, details string) {
	response := APIResponse{
		Success: false,
		Error: &APIError{
			Code:    code,
			Message: message,
			Details: details,
		},
	}
	c.JSON(statusCode, response)
}

// Helper functions for common error responses

// BadRequest returns 400 Bad Request
func BadRequest(c *gin.Context, message string, details ...string) {
	detail := ""
	if len(details) > 0 {
		detail = details[0]
	}
	ErrorResponse(c, http.StatusBadRequest, ErrCodeBadRequest, message, detail)
}

// Unauthorized returns 401 Unauthorized
func Unauthorized(c *gin.Context, message string) {
	if message == "" {
		message = "Unauthorized"
	}
	ErrorResponse(c, http.StatusUnauthorized, ErrCodeUnauthorized, message, "")
}

// Forbidden returns 403 Forbidden
func Forbidden(c *gin.Context, message string) {
	if message == "" {
		message = "Forbidden"
	}
	ErrorResponse(c, http.StatusForbidden, ErrCodeForbidden, message, "")
}

// NotFound returns 404 Not Found
func NotFound(c *gin.Context, message string) {
	if message == "" {
		message = "Resource not found"
	}
	ErrorResponse(c, http.StatusNotFound, ErrCodeNotFound, message, "")
}

// InternalError returns 500 Internal Server Error
func InternalError(c *gin.Context, message string, details ...string) {
	detail := ""
	if len(details) > 0 {
		detail = details[0]
	}
	if message == "" {
		message = "Internal server error"
	}
	ErrorResponse(c, http.StatusInternalServerError, ErrCodeInternal, message, detail)
}

// ValidationError returns 400 Bad Request for validation errors
func ValidationError(c *gin.Context, message string, details ...string) {
	detail := ""
	if len(details) > 0 {
		detail = details[0]
	}
	ErrorResponse(c, http.StatusBadRequest, ErrCodeValidation, message, detail)
}

// Conflict returns 409 Conflict
func Conflict(c *gin.Context, message string) {
	if message == "" {
		message = "Resource conflict"
	}
	ErrorResponse(c, http.StatusConflict, ErrCodeConflict, message, "")
}

// OK returns 200 OK with data
func OK(c *gin.Context, data interface{}, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	SuccessResponse(c, http.StatusOK, data, msg)
}

// Created returns 201 Created with data
func Created(c *gin.Context, data interface{}, message ...string) {
	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}
	SuccessResponse(c, http.StatusCreated, data, msg)
}

// TooManyRequests returns 429 Too Many Requests
func TooManyRequests(c *gin.Context, message string) {
	if message == "" {
		message = "Too many requests"
	}
	ErrorResponse(c, http.StatusTooManyRequests, ErrCodeRateLimitExceeded, message, "")
}

