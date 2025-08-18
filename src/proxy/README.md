# Nginx Proxy Service

This is a custom Nginx proxy service for the Simple Chat application.

## Features

- **Reverse Proxy**: Routes requests to appropriate backend services
- **Load Balancing**: Distributes traffic across multiple upstream servers
- **Rate Limiting**: Protects against abuse with configurable rate limits
- **WebSocket Support**: Handles WebSocket connections for real-time chat
- **Security Headers**: Adds security headers to all responses
- **Gzip Compression**: Compresses responses for better performance
- **Health Checks**: Built-in health check endpoint
- **Monitoring**: Nginx status endpoint for monitoring

## Architecture

```
Client Request → Nginx Proxy → Backend Services
                ├── / → Frontend (React)
                ├── /api/ → Backend API
                ├── /ws/ → WebSocket (Backend)
                └── /socket/ → Socket Service
```

## Configuration

### Main Configuration (`nginx.conf`)
- Worker processes and connections
- Logging format
- Gzip compression settings
- Rate limiting zones
- Upstream server definitions

### Server Configuration (`conf.d/default.conf`)
- Request routing rules
- Proxy settings
- Security headers
- Rate limiting rules
- WebSocket handling

## Routes

| Path | Service | Description |
|------|---------|-------------|
| `/` | Frontend | React application |
| `/api/` | Backend | REST API endpoints |
| `/ws/` | Backend | WebSocket connections |
| `/socket/` | Socket Service | Dedicated WebSocket service |
| `/health` | Proxy | Health check endpoint |
| `/nginx_status` | Proxy | Nginx status monitoring |

## Rate Limiting

- **API**: 10 requests/second with burst of 20
- **WebSocket**: 30 requests/second with burst of 50
- **Login**: 5 requests/minute

## Development

### Building the Image
```bash
# Production
docker build -t simple-chat-proxy .

# Development
docker build -f Dockerfile.dev -t simple-chat-proxy-dev .
```

### Running Locally
```bash
# Production
docker run -p 80:80 simple-chat-proxy

# Development with volume mounts
docker run -p 80:80 -v $(pwd)/nginx.conf:/etc/nginx/nginx.conf -v $(pwd)/conf.d:/etc/nginx/conf.d simple-chat-proxy-dev
```

### Testing Configuration
```bash
# Test nginx configuration
docker exec simple-chat-proxy nginx -t

# Check nginx status
curl http://localhost/nginx_status

# Health check
curl http://localhost/health
```

## Environment Variables

- `NGINX_PORT`: Port for HTTP (default: 80)
- `NGINX_SSL_PORT`: Port for HTTPS (default: 443)

## Monitoring

### Health Check
```bash
curl http://localhost/health
# Returns: "healthy"
```

### Nginx Status
```bash
curl http://localhost/nginx_status
# Returns nginx status information
```

### Logs
```bash
# Access logs
docker logs simple-chat-proxy

# Error logs
docker exec simple-chat-proxy tail -f /var/log/nginx/error.log
```

## Security

### Security Headers
- `X-Frame-Options`: Prevents clickjacking
- `X-XSS-Protection`: XSS protection
- `X-Content-Type-Options`: Prevents MIME type sniffing
- `Referrer-Policy`: Controls referrer information
- `Content-Security-Policy`: Content security policy

### Rate Limiting
- API endpoints are rate limited to prevent abuse
- WebSocket connections have higher limits for real-time communication
- Login endpoints have strict rate limiting

## Performance

### Gzip Compression
- Compresses text-based files (HTML, CSS, JS, JSON)
- Reduces bandwidth usage
- Improves page load times

### Caching
- Static files are cached for 1 year
- Proper cache headers for optimal performance

### Connection Handling
- Keep-alive connections for better performance
- Proper timeout settings for WebSocket connections
