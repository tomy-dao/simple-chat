# Backend Architecture Documentation

## Tổng quan

Backend được xây dựng theo kiến trúc **Clean Architecture** với các layer rõ ràng, hỗ trợ logging, tracing và job processing.

## Kiến trúc Layers

### 1. Transport Layer (`transport/http/`)
**Chức năng**: Xử lý HTTP requests/responses, routing, middleware

**Cấu trúc**:
- `transport.go`: Khởi tạo Gin engine và setup middleware
- `router.go`: Định nghĩa routes và mapping với handlers
- `handler.go`: HTTP handlers, chuyển đổi HTTP requests thành endpoint calls
- `middleware.go`: Authentication, CORS, JSON content-type, OpenTelemetry tracing
- `response.go`: Response helpers

**Flow**:
```
HTTP Request → Middleware → Router → Handler → Endpoint
```

**Middleware Stack**:
1. OpenTelemetry Tracing (tạo span cho mỗi request)
2. CORS
3. JSON Content-Type
4. Authentication (cho protected routes)

### 2. Endpoint Layer (`endpoint/`)
**Chức năng**: Interface giữa transport và service layer, xử lý request/response transformation

**Cấu trúc**:
- `auth.go`: Auth endpoints (Register, Login, Logout, GetMe, GetUsers)
- `conversation.go`: Conversation endpoints
- `message.go`: Message endpoints
- `initial.go`: Khởi tạo tất cả endpoints

**Pattern**:
- Endpoints nhận `*model.RequestContext` và request data
- Gọi service layer
- Trả về `model.Response[T]`

### 3. Service Layer (`service/`)
**Chức năng**: Business logic, validation, orchestration

**Cấu trúc**:
- `auth/service.go`: Authentication, JWT token management, user management
- `conversation/service.go`: Conversation creation, retrieval
- `message/service.go`: Message creation, retrieval, broadcasting
- `initial/service.go`: Khởi tạo tất cả services
- `common/model.go`: Common parameters (Repo, Client)

**Dependencies**:
- Services phụ thuộc vào Repository (data access)
- Services có thể phụ thuộc vào các services khác (MessageService → ConversationService)

**Pattern**:
- Mỗi service implement một interface
- Services nhận `*model.RequestContext` để logging và tracing
- Trả về `model.Response[T]` với standardized error handling

### 4. Repository Layer (`infra/repo/`)
**Chức năng**: Data access, database operations

**Cấu trúc**:
- `repo.go`: Repository interface và factory
- `user_repo.go`: User CRUD operations
- `conversation_repo.go`: Conversation CRUD operations
- `participant_repo.go`: Participant management
- `message_repo.go`: Message CRUD operations

**Pattern**:
- Repository pattern với interface
- Sử dụng GORM cho database operations
- Auto-migration cho models
- Trả về `model.Response[T]` cho consistency

**Database**:
- MySQL với GORM ORM
- Auto-migration on startup
- Connection pooling

## Logging & Tracing

### Logging (`util/logger/logger.go`)

**Framework**: Zerolog

**Features**:
- Structured logging với JSON format (production) hoặc Console format (development)
- Log levels: Debug, Info, Warn, Error, Fatal
- Context-aware logging với RequestContext
- Tự động thêm trace_id, span_id, user_id, session_id vào logs

**Usage**:
```go
logger.Info(reqCtx, "Message", map[string]interface{}{
    "key": "value",
})
logger.Error(reqCtx, "Error message", err, map[string]interface{}{})
```

**Log Format**:
- Development: Console format (human-readable)
- Production/Docker: JSON format (cho Loki)
- Environment variable: `LOG_FORMAT=json|console`

### Tracing (`util/logger/tracer.go`)

**Framework**: OpenTelemetry

**Features**:
- Distributed tracing với OpenTelemetry
- Export traces đến Jaeger qua OTLP HTTP
- Tự động tạo spans cho HTTP requests
- Span events từ logs

**Configuration**:
- `OTEL_SERVICE_NAME`: Service name (default: từ InitTracer parameter)
- `OTEL_EXPORTER_OTLP_ENDPOINT`: Jaeger endpoint (default: http://localhost:4318)

**Integration**:
- Middleware tự động tạo span cho mỗi HTTP request
- Logger tự động thêm log events vào spans
- RequestContext chứa span để propagate qua layers

**Observability Stack**:
- **Loki**: Log aggregation (từ Promtail)
- **Grafana**: Visualization và dashboards
- **Jaeger**: Distributed tracing
- **Promtail**: Log shipper (Docker logs → Loki)

## Job System (`cmd/job/`)

**Framework**: Cobra CLI

**Cấu trúc**:
- `job.go`: Job command definition và execution
- Jobs chạy như standalone processes

**Job Types**:
1. **cleanup**: Cleanup job (ví dụ: xóa old data)
2. **sync**: Sync job (ví dụ: sync data với external systems)

**Usage**:
```bash
./simple-chat job --type cleanup
./simple-chat job --type sync
```

**Features**:
- Tự động khởi tạo tracer cho jobs
- Structured logging với trace context
- Graceful shutdown

## Models & Response Pattern

### RequestContext (`model/request_context.go`)

**Chức năng**: Context wrapper với type-safe accessors

**Fields**:
- `ctx`: Underlying context.Context
- `Token`: JWT token
- `UserID`: Authenticated user ID
- `SessionID`: Session ID từ JWT
- `span`: OpenTelemetry span

**Methods**:
- `NewRequestContext(ctx)`: Tạo từ context
- `WithToken/WithUserID/WithSessionID/WithClaims`: Immutable updates
- `TraceID()/SpanID()`: Lấy trace/span IDs
- `Span()`: Lấy OpenTelemetry span

### Response Pattern (`model/response.go`)

**Generic Response Type**:
```go
type Response[T any] struct {
    Code    int     `json:"code"`
    Message string  `json:"message"`
    Data    T       `json:"data"`
    Errors  []Error `json:"errors,omitempty"`
}
```

**Response Codes**:
- `200`: Success
- `400`: Bad Request
- `401`: Unauthorized
- `403`: Forbidden
- `404`: Not Found
- `409`: Conflict
- `422`: Validation Error
- `500`: Internal Server Error

**Helper Functions**:
- `SuccessResponse(data, message)`: Tạo success response
- `BadRequest/Unauthorized/NotFound/Conflict/InternalError(message)`: Tạo error responses
- `ValidationError/ValidationErrorWithErrors`: Validation errors

## Configuration (`config/config.go`)

**Environment Variables**:
- `HTTP_PORT`: HTTP server port (default: 80)
- `HOST`: Server host (default: 0.0.0.0)
- `DB_HOST`: Database host (default: localhost)
- `DB_PORT`: Database port (default: 3306)
- `DB_USER`: Database user (default: root)
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name (default: simple_chat)
- `JWT_SECRET`: JWT signing secret
- `SOCKET_SERVER_URL`: Socket server URL
- `SOCKET_TOKEN`: Socket server token
- `OTEL_SERVICE_NAME`: OpenTelemetry service name
- `OTEL_EXPORTER_OTLP_ENDPOINT`: Jaeger endpoint
- `LOG_FORMAT`: Log format (json|console)
- `ENV`: Environment (development|production|docker)

**Load Config**:
- `config.LoadConfig()`: Load từ environment variables
- Fallback values cho tất cả configs

## Command Structure (`cmd/`)

**Cobra CLI Commands**:
- `server`: Start HTTP server
- `job`: Run background jobs

**Entry Point**:
- `main.go`: Khởi tạo và execute root command
- `cmd/root.go`: Root command definition
- `cmd/server.go`: Server command
- `cmd/job/job.go`: Job command

## Testing Structure (`test/`)

**Organization**:
- `test/service/`: Service layer unit tests
- `test/integration/`: Integration tests
- `test/infra/repo/`: Repository tests
- `test/mocks/`: Mock implementations

**Test Helpers**:
- `service/auth/test_helper.go`: Auth service test helpers
- `integration/helper.go`: Integration test helpers

## Data Flow

### Request Flow:
```
HTTP Request
  ↓
Transport Layer (middleware, routing)
  ↓
Handler (HTTP → Endpoint transformation)
  ↓
Endpoint Layer (request validation)
  ↓
Service Layer (business logic)
  ↓
Repository Layer (database operations)
  ↓
Database
```

### Response Flow:
```
Database
  ↓
Repository (model.Response[T])
  ↓
Service (business logic, validation)
  ↓
Endpoint (response transformation)
  ↓
Handler (JSON serialization)
  ↓
HTTP Response
```

## Error Handling

**Pattern**:
- Tất cả layers trả về `model.Response[T]`
- Error codes standardized
- Error messages user-friendly
- Errors array cho validation errors

**Propagation**:
- Errors được propagate từ repository → service → endpoint → handler
- Không throw exceptions, chỉ return error responses

## Security

**Authentication**:
- JWT tokens với HS256 signing
- Token validation trong middleware
- User context trong RequestContext

**Password Security**:
- bcrypt hashing với DefaultCost
- Passwords không bao giờ trả về trong responses

## Socket Integration

**Client** (`client/`):
- Socket client để broadcast messages
- Integration với socket server
- Broadcast events khi có message mới

## Swagger Documentation

**Location**: `docs/`

**Generation**:
```bash
make swagger
# hoặc
swag init -g main.go -o ./docs --parseDependency --parseInternal --parseDepth 2
```

**Access**: `http://localhost:80/swagger/index.html`

**Annotations**:
- Swagger annotations trong `handler.go`
- Generic types: `model.Response[T]` → `model.Response-T` trong docs

## Best Practices

1. **Context Propagation**: Luôn pass `RequestContext` qua tất cả layers
2. **Error Handling**: Sử dụng `model.Response[T]` cho tất cả responses
3. **Logging**: Log ở mỗi layer với context information
4. **Tracing**: Spans tự động tạo, logs tự động thêm vào spans
5. **Testing**: Unit tests cho services, integration tests cho endpoints
6. **Type Safety**: Sử dụng generic types cho type safety
7. **Dependency Injection**: Services nhận dependencies qua constructor

## File Structure

```
backend/
├── cmd/              # CLI commands (server, job)
├── config/           # Configuration
├── endpoint/         # Endpoint layer
├── service/          # Service layer (business logic)
│   ├── auth/
│   ├── conversation/
│   ├── message/
│   ├── initial/
│   └── common/
├── infra/repo/       # Repository layer (data access)
├── model/            # Domain models
├── transport/http/   # HTTP transport layer
├── util/logger/     # Logging và tracing
├── client/          # Socket client
├── docs/            # Swagger documentation
└── test/            # Tests
```

## Dependencies

**Main Libraries**:
- `gin-gonic/gin`: HTTP framework
- `gorm.io/gorm`: ORM
- `golang-jwt/jwt/v5`: JWT tokens
- `rs/zerolog`: Structured logging
- `opentelemetry.io/otel`: Distributed tracing
- `spf13/cobra`: CLI framework
- `swaggo/swag`: Swagger documentation

## Environment Setup

**Development**:
- Console logging format
- Local database
- Local Jaeger

**Production/Docker**:
- JSON logging format (cho Loki)
- Docker database
- Jaeger trong Docker network

## Monitoring & Observability

**Metrics**:
- Logs → Loki (via Promtail)
- Traces → Jaeger
- Dashboards → Grafana

**Log Fields**:
- `trace_id`: Distributed trace ID
- `span_id`: Current span ID
- `user_id`: Authenticated user ID
- `session_id`: Session ID
- Custom fields từ business logic

