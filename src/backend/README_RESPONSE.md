# Standard Response & Error Handling

Tài liệu về chuẩn response và error handling cho Simple Chat API.

## Response Structure

Tất cả API responses đều tuân theo cấu trúc chuẩn:

```json
{
  "success": true,
  "data": { ... },
  "message": "Optional success message",
  "error": null
}
```

### Success Response

```json
{
  "success": true,
  "data": {
    "id": 1,
    "username": "john_doe"
  },
  "message": "User registered successfully"
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request body",
    "details": "username is required"
  }
}
```

## Error Codes

| Code | HTTP Status | Mô tả |
|------|-------------|-------|
| `VALIDATION_ERROR` | 400 | Lỗi validation input |
| `BAD_REQUEST` | 400 | Request không hợp lệ |
| `UNAUTHORIZED` | 401 | Chưa đăng nhập hoặc token không hợp lệ |
| `FORBIDDEN` | 403 | Không có quyền truy cập |
| `NOT_FOUND` | 404 | Resource không tồn tại |
| `CONFLICT` | 409 | Resource đã tồn tại (conflict) |
| `INTERNAL_ERROR` | 500 | Lỗi server |

## Helper Functions

### Success Responses

```go
// 200 OK
OK(c, data, "Optional message")

// 201 Created
Created(c, data, "Optional message")
```

### Error Responses

```go
// 400 Bad Request
BadRequest(c, "Error message", "Optional details")

// 400 Validation Error
ValidationError(c, "Validation error", "Field details")

// 401 Unauthorized
Unauthorized(c, "Unauthorized message")

// 403 Forbidden
Forbidden(c, "Forbidden message")

// 404 Not Found
NotFound(c, "Not found message")

// 409 Conflict
Conflict(c, "Conflict message")

// 500 Internal Error
InternalError(c, "Internal error", "Optional details")
```

## Usage Examples

### Success Response

```go
func (h *handler) GetMe() gin.HandlerFunc {
    return func(c *gin.Context) {
        user, err := h.endpoints.Auth.GetMe(c.Request.Context())
        if err != nil {
            BadRequest(c, err.Error())
            return
        }
        OK(c, user)
    }
}
```

### Error Response

```go
func (h *handler) Register() gin.HandlerFunc {
    return func(c *gin.Context) {
        var req endpoint.RegisterRequest
        if err := c.ShouldBindJSON(&req); err != nil {
            ValidationError(c, "Invalid request body", err.Error())
            return
        }
        // ...
    }
}
```

### Created Response

```go
func (h *handler) CreateConversation() gin.HandlerFunc {
    return func(c *gin.Context) {
        // ... create conversation
        Created(c, conversation, "Conversation created successfully")
    }
}
```

## Response Structure Details

### APIResponse

```go
type APIResponse struct {
    Success bool        `json:"success"`  // true nếu thành công, false nếu có lỗi
    Data    interface{} `json:"data,omitempty"`  // Data trả về (chỉ có khi success)
    Error   *APIError   `json:"error,omitempty"`  // Error object (chỉ có khi có lỗi)
    Message string      `json:"message,omitempty"`  // Optional message
}
```

### APIError

```go
type APIError struct {
    Code    string `json:"code"`     // Error code (VALIDATION_ERROR, etc.)
    Message string `json:"message"` // Error message
    Details string `json:"details,omitempty"`  // Optional details
}
```

## Best Practices

1. **Luôn sử dụng helper functions** thay vì `c.JSON()` trực tiếp
2. **Sử dụng đúng error code** cho từng loại lỗi
3. **Thêm message** cho success responses khi có thể
4. **Thêm details** cho validation errors để client biết field nào lỗi
5. **Consistent format** - tất cả responses đều theo cùng một format

## Migration

Tất cả handlers đã được update để sử dụng chuẩn response mới:
- ✅ `GetMe`
- ✅ `Register`
- ✅ `Login`
- ✅ `Logout`
- ✅ `CreateConversation`
- ✅ `GetConversationByUserID`
- ✅ `GetConversations`
- ✅ `CreateMessage`
- ✅ `GetMessagesByConversationID`
- ✅ `GetUsers`

