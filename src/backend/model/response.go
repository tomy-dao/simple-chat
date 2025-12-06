package model

// Response code constants
const (
	CodeSuccess        = 200 // Success
	CodeBadRequest     = 400 // Bad Request
	CodeUnauthorized   = 401 // Unauthorized
	CodeForbidden      = 403 // Forbidden
	CodeNotFound       = 404 // Not Found
	CodeConflict       = 409 // Conflict
	CodeValidation     = 422 // Validation Error
	CodeInternalError  = 500 // Internal Server Error
)

// CodeMessages maps response codes to their default messages; when adding a new code,
// remember to update this map to keep error responses consistent.
var CodeMessages = map[int]string{
	CodeSuccess:       "Success",
	CodeBadRequest:    "Bad Request",
	CodeUnauthorized:  "Unauthorized",
	CodeForbidden:     "Forbidden",
	CodeNotFound:      "Not Found",
	CodeConflict:      "Conflict",
	CodeValidation:    "Validation Error",
	CodeInternalError: "Internal Server Error",
}

// GetCodeMessage returns the message string for a given code
// Returns the message from CodeMessages map, or empty string if not found
func GetCodeMessage(code int) string {
	if msg, ok := CodeMessages[code]; ok {
		return msg
	}
	return ""
}



// Error represents an error with code and message
// swagger:model Error
type Error struct {
	Code    int    `json:"code" example:"500"`
	Message string `json:"message" example:"Internal Server Error"`
}

// Response represents a standard API response
// Data is generic type T
type Response[T any] struct {
	Code    int     `json:"code"`
	Message string  `json:"message"`
	Data    T       `json:"data"`
	Errors  []Error `json:"errors,omitempty"` // Array of errors with code and message
}

// OK checks if the response is successful
// Returns true if Code is 200 (SUCCESS), false otherwise
func (r *Response[T]) OK() bool {
	return r.Code == CodeSuccess
}

// GetErrorString returns the error message string from the response
// Returns Message if not success, otherwise returns empty string
func (r *Response[T]) ErrorString() string {
	if !r.OK() {
		return r.Message
	}
	return ""
}

func (r *Response[T]) ErrorCodeMessage() string {
	if !r.OK() {
		return GetCodeMessage(r.Code)
	}
	return ""
}

// ErrorResponse creates an error response with custom code and message
func ErrorResponse[T any](code int, message string) Response[T] {
	var zero T
	return Response[T]{
		Code:    code,
		Message: message,
		Data:    zero,
		Errors:  []Error{{Code: code, Message: message}},
	}
}

// ErrorArray creates an error response with array of errors
func ErrorArray[T any](code int, message string, responseErrors []Error) Response[T] {
	var zero T
	return Response[T]{
		Code:    code,
		Message: message,
		Data:    zero,
		Errors:  responseErrors,
	}
}

// BadRequest creates a 400 Bad Request error response
func BadRequest[T any](message string) Response[T] {
	return ErrorResponse[T](CodeBadRequest, message)
}

// BadRequestWithErrors creates a 400 Bad Request error response with array of errors
func BadRequestWithErrors[T any](message string, errors []Error) Response[T] {
	return ErrorArray[T](CodeBadRequest, message, errors)
}

// Unauthorized creates a 401 Unauthorized error response
func Unauthorized[T any](message string) Response[T] {
	return ErrorResponse[T](CodeUnauthorized, message)
}

// Forbidden creates a 403 Forbidden error response
func Forbidden[T any](message string) Response[T] {
	return ErrorResponse[T](CodeForbidden, message)
}

// NotFound creates a 404 Not Found error response
func NotFound[T any](message string) Response[T] {
	return ErrorResponse[T](CodeNotFound, message)
}

// Conflict creates a 409 Conflict error response
func Conflict[T any](message string) Response[T] {
	return ErrorResponse[T](CodeConflict, message)
}

// ValidationError creates a 422 Validation Error response
func ValidationError[T any](message string) Response[T] {
	return ErrorResponse[T](CodeValidation, message)
}

// ValidationErrorWithErrors creates a 422 Validation Error response with array of errors
func ValidationErrorWithErrors[T any](message string, errors []Error) Response[T] {
	return ErrorArray[T](CodeValidation, message, errors)
}

// InternalError creates a 500 Internal Server Error response
func InternalError[T any](message string) Response[T] {
	return ErrorResponse[T](CodeInternalError, message)
}

// SuccessResponse creates a success response with data and message
func SuccessResponse[T any](data T, message string) Response[T] {
	return Response[T]{
		Code:    CodeSuccess,
		Message: message,
		Data:    data,
	}
}


