# Go-Chi Router Implementation

This project has been converted to use the go-chi router instead of the standard `http.ServeMux`. Go-chi provides a lightweight, expressive HTTP router with middleware support.

## Features

### 1. **Router Structure**
- **Chi Router**: Lightweight, fast HTTP router
- **Middleware Support**: Built-in and custom middleware
- **Route Groups**: Organized API structure
- **CORS Support**: Cross-origin resource sharing enabled
- **Authentication**: JWT-based authentication system

### 2. **API Endpoints**

#### Base URL: `http://localhost:80`

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | API information and available endpoints |
| GET | `/health` | Health check endpoint |
| POST | `/api/v1/auth/register` | Register new user |
| POST | `/api/v1/auth/login` | User login |
| POST | `/api/v1/auth/logout` | User logout |
| GET | `/api/v1/users` | Get all users |
| POST | `/api/v1/utils/uppercase` | Convert string to uppercase |

### 3. **Middleware Stack**

The application includes the following middleware:

- **RequestID**: Adds unique request ID to each request
- **RealIP**: Extracts real IP from headers
- **Logger**: Logs HTTP requests
- **Recoverer**: Recovers from panics
- **CleanPath**: Cleans URL paths
- **GetHead**: Handles HEAD requests
- **Timeout**: 60-second request timeout
- **CORS**: Cross-origin resource sharing

### 4. **Project Structure**

```
transport/http/
├── transport.go      # Main transport setup
├── router.go         # Route definitions
├── handler.go        # HTTP handlers
└── middleware.go     # Middleware functions

service/auth/
└── service.go        # Authentication service

endpoint/
├── auth.go           # Authentication endpoints
└── user.go           # User endpoints
```

## API Usage Examples

### 1. **Get API Information**
```bash
curl http://localhost/
```

Response:
```json
{
  "message": "Local Service API",
  "version": "1.0.0",
  "endpoints": {
    "health": "/health",
    "register": "/api/v1/auth/register",
    "login": "/api/v1/auth/login",
    "logout": "/api/v1/auth/logout",
    "users": "/api/v1/users",
    "uppercase": "/api/v1/utils/uppercase"
  }
}
```

### 2. **Health Check**
```bash
curl http://localhost/health
```

Response:
```json
{
  "status": "OK",
  "timestamp": "2024-01-01T12:00:00Z",
  "service": "local-service"
}
```

### 3. **User Registration**
```bash
curl -X POST http://localhost/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"user_name": "john_doe", "password": "secure_password"}'
```

Response:
```json
{
  "id": 1,
  "user_name": "john_doe",
  "password": "",
  "created_at": "2024-01-01T12:00:00Z",
  "updated_at": "2024-01-01T12:00:00Z"
}
```

### 4. **User Login**
```bash
curl -X POST http://localhost/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"user_name": "john_doe", "password": "secure_password"}'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 5. **User Logout**
```bash
curl -X POST http://localhost/api/v1/auth/logout \
  -H "Content-Type: application/json" \
  -d '{"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."}'
```

Response:
```json
{
  "message": "Logged out successfully"
}
```

### 6. **Get Users**
```bash
curl http://localhost/api/v1/users
```

Response:
```json
{
  "users": [
    {
      "id": 1,
      "user_name": "john_doe",
      "password": "",
      "created_at": "2024-01-01T12:00:00Z",
      "updated_at": "2024-01-01T12:00:00Z"
    }
  ],
  "err": ""
}
```

### 7. **Uppercase String**
```bash
curl -X POST http://localhost/api/v1/utils/uppercase \
  -H "Content-Type: application/json" \
  -d '{"s": "hello world"}'
```

Response:
```json
{
  "v": "HELLO WORLD",
  "err": ""
}
```

## Key Changes Made

### 1. **Dependencies Added**
- `github.com/go-chi/chi/v5`: Main router
- `github.com/go-chi/cors`: CORS middleware
- `github.com/golang-jwt/jwt/v5`: JWT authentication
- `github.com/google/uuid`: UUID generation
- `crypto/bcrypt`: Password hashing

### 2. **Transport Layer**
- Replaced `http.ServeMux` with `chi.Router`
- Added comprehensive middleware stack
- Organized route structure with API versioning
- Integrated authentication service

### 3. **Authentication System**
- JWT-based authentication
- Password hashing with bcrypt
- User registration and login
- Session management

### 4. **Handler Layer**
- Converted from go-kit transport servers to `http.HandlerFunc`
- Direct endpoint integration
- Proper error handling and JSON responses
- Authentication handlers

### 5. **Middleware Organization**
- Centralized middleware configuration
- Custom middleware functions
- CORS support for cross-origin requests

## Benefits of Go-Chi

1. **Performance**: Lightweight and fast routing
2. **Middleware**: Rich middleware ecosystem
3. **Expressiveness**: Clean, readable route definitions
4. **Flexibility**: Easy to extend and customize
5. **Standards**: Follows Go HTTP standards

## Running the Application

1. **Start the application**:
   ```bash
   go run main.go
   ```

2. **Test endpoints**:
   ```bash
   # Test health check
   curl http://localhost/health
   
   # Test API info
   curl http://localhost/
   
   # Test user registration
   curl -X POST http://localhost/api/v1/auth/register \
     -H "Content-Type: application/json" \
     -d '{"user_name": "test_user", "password": "password123"}'
   
   # Test user login
   curl -X POST http://localhost/api/v1/auth/login \
     -H "Content-Type: application/json" \
     -d '{"user_name": "test_user", "password": "password123"}'
   ```

## Security Considerations

1. **JWT Secret**: In production, use environment variables for JWT secret
2. **Password Hashing**: Passwords are hashed using bcrypt
3. **CORS**: Configure CORS properly for production
4. **Token Expiration**: JWT tokens expire after 24 hours
5. **Input Validation**: Add proper input validation for production

## Adding New Endpoints

To add new endpoints, follow this pattern:

1. **Add to router.go**:
   ```go
   r.Route("/api/v1/users", func(r chi.Router) {
       r.Get("/", h.GetUsers())
       r.Post("/", h.CreateUser())        // New endpoint
       r.Get("/{id}", h.GetUser())        // New endpoint
   })
   ```

2. **Add to handler.go**:
   ```go
   func (h *handler) CreateUser() http.HandlerFunc {
       return func(w http.ResponseWriter, r *http.Request) {
           // Implementation
       }
   }
   ```

3. **Add to service.go**:
   ```go
   func (svc *service) CreateUser(ctx context.Context, user *model.User) (*model.User, error) {
       // Implementation
   }
   ```

## Middleware Customization

You can customize middleware in `middleware.go`:

```go
// Add custom middleware
func CustomMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Custom logic
        next.ServeHTTP(w, r)
    })
}
```

Then add it to `SetupMiddleware()`:

```go
func SetupMiddleware(r chi.Router) {
    // ... existing middleware
    r.Use(CustomMiddleware)
}
```
