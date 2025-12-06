# Swagger Documentation

Swagger (OpenAPI) documentation cho Simple Chat API.

## Setup

1. **Install Swagger CLI**:
   ```bash
   go install github.com/swaggo/swag/cmd/swag@latest
   ```

2. **Generate Swagger Documentation**:
   ```bash
   make swagger
   ```

## Swagger Annotations

Swagger annotations được đặt trực tiếp trong `transport/http/handler.go` cho mỗi handler function.

### Format chuẩn:

```go
// HandlerName godoc
// @Summary Brief description
// @Tags tag1
// @Security BearerAuth (nếu cần auth)
// @Param request body endpoint.RequestType true "Request description"
// @Param paramName path int true "Path parameter"
// @Success 200 {object} model.ResponseType
// @Router /api/v1/endpoint [method]
func (h *handler) HandlerName() gin.HandlerFunc {
    // Handler implementation
}
```

### Ví dụ:

```go
// Register godoc
// @Summary Register a new user
// @Tags auth
// @Param request body endpoint.RegisterRequest true "Register Request"
// @Success 200 {object} model.User
// @Router /register [post]
func (h *handler) Register() gin.HandlerFunc {
    // ...
}
```

## Swagger Annotations Reference

- `@Summary` - Mô tả ngắn gọn của endpoint
- `@Tags` - Nhóm endpoints (auth, conversations, messages, users)
- `@Security BearerAuth` - Yêu cầu authentication
- `@Param` - Request parameters (body, path, query)
- `@Success` - Response type và status code
- `@Router` - Route path và HTTP method

## Access Swagger UI

Khi server đang chạy:
- **Swagger UI**: http://localhost:80/swagger/index.html
- **API Base**: http://localhost:80/api/v1

## Updating Documentation

1. Thêm/sửa handler trong `transport/http/handler.go`
2. Thêm Swagger annotations cho handler mới
3. Chạy `make swagger` để regenerate documentation
4. Restart server để xem thay đổi

## Main API Info

API info được định nghĩa trong `main.go`:
- Title, version, description
- Base path, host
- Security definitions (Bearer Auth)
