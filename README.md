# Simple Chat Application

A real-time chat application built with Go backend, React frontend, and MySQL database using Docker.

## Architecture

- **Backend**: Go with Gin framework
- **Frontend**: React with TypeScript
- **Socket Service**: Go WebSocket service for real-time communication
- **Database**: MySQL
- **Cache**: Redis for socket clustering
- **Proxy**: Nginx reverse proxy
- **Containerization**: Docker & Docker Compose

## Prerequisites

- Docker
- Docker Compose

## Quick Start

### 1. Clone the repository
```bash
git clone <repository-url>
cd simple-chat
```

### 2. Set up environment variables
```bash
# Copy the example environment file
cp env.example .env

# Edit the .env file with your configuration
nano .env
```

**Important**: All environment variables are now required. Make sure to set all values in your `.env` file before running the application.

### 3. Run the application

#### Production mode:
```bash
docker-compose up -d
```

#### Development mode (with hot reloading):
```bash
docker-compose -f docker-compose.dev.yaml up -d
```

### 4. Access the application
- Frontend: http://localhost
- Backend API: http://localhost/api
- WebSocket: ws://localhost/ws
- Socket Service: ws://localhost/socket
- MySQL: localhost:3306
- Redis: localhost:6379

## Environment Variables

### MySQL Configuration
- `MYSQL_DATABASE`: Database name (default: simple_chat)
- `MYSQL_USER`: Database user (default: chat_user)
- `MYSQL_PASSWORD`: Database password (default: password)
- `MYSQL_ROOT_PASSWORD`: MySQL root password (default: root123)
- `MYSQL_PORT`: Database port (default: 3306)

### Backend Configuration
- `BACKEND_PORT`: Backend server port (default: 8080)
- `JWT_SECRET`: JWT secret key for authentication
- `JWT_EXPIRY`: JWT token expiry time (default: 24h)
- `CORS_ORIGIN`: CORS allowed origin (default: http://localhost:3000)
- `ENV`: Environment (development/production)

### Nginx Proxy Configuration
- `NGINX_PORT`: Nginx proxy port (default: 80)
- `NGINX_SSL_PORT`: Nginx SSL port (default: 443)

### Socket Service Configuration
- `SOCKET_PORT`: Socket service port (default: 8081)
- `REDIS_URL`: Redis connection URL for socket clustering

### Frontend Configuration
- `FRONTEND_PORT`: Frontend server port (default: 3000)
- `REACT_APP_API_URL`: Backend API URL
- `REACT_APP_WS_URL`: WebSocket URL for real-time chat
- `REACT_APP_SOCKET_URL`: Socket service URL
- `REACT_APP_NAME`: Application name
- `REACT_APP_VERSION`: Application version
- `NODE_ENV`: Node.js environment

## Development

### Backend Development
The Go backend uses Air for hot reloading in development mode. Any changes to Go files will automatically restart the server.

### Frontend Development
The React frontend uses the standard development server with hot reloading. Any changes to React files will automatically refresh the browser.

### Database
The MySQL database is initialized with sample data from `init.sql`. The database data is persisted in Docker volumes.

## Docker Commands

### Start services
```bash
# Production
docker-compose up -d

# Development
docker-compose -f docker-compose.dev.yaml up -d
```

### Stop services
```bash
# Production
docker-compose down

# Development
docker-compose -f docker-compose.dev.yaml down
```

### View logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f backend
docker-compose logs -f frontend
docker-compose logs -f postgres
```

### Rebuild services
```bash
# Production
docker-compose up -d --build

# Development
docker-compose -f docker-compose.dev.yaml up -d --build
```

### Clean up
```bash
# Stop and remove containers, networks
docker-compose down

# Stop and remove containers, networks, volumes
docker-compose down -v

# Remove all unused containers, networks, images
docker system prune -a
```

## Project Structure

```
simple-chat/
├── docker-compose.yaml          # Production Docker Compose
├── docker-compose.dev.yaml      # Development Docker Compose
├── env.example                  # Environment variables template
├── init.sql                     # Database initialization script
├── README.md                    # This file
└── src/
    ├── backend/
    │   ├── Dockerfile           # Production backend Dockerfile
    │   ├── Dockerfile.dev       # Development backend Dockerfile
    │   └── .air.toml           # Air configuration for hot reloading
    ├── socket/
    │   ├── Dockerfile           # Production socket Dockerfile
    │   ├── Dockerfile.dev       # Development socket Dockerfile
    │   └── .air.toml           # Air configuration for hot reloading
    ├── proxy/
    │   ├── Dockerfile           # Production proxy Dockerfile
    │   ├── Dockerfile.dev       # Development proxy Dockerfile
    │   ├── nginx.conf          # Main nginx configuration
    │   ├── conf.d/
    │   │   └── default.conf    # Proxy routing rules
    │   ├── start.sh            # Development startup script
    │   └── README.md           # Proxy service documentation
    └── frontend/
        ├── Dockerfile           # Production frontend Dockerfile
        ├── Dockerfile.dev       # Development frontend Dockerfile
        └── nginx.conf          # Frontend nginx configuration
```

## Features

- Real-time chat with WebSocket support
- Separate WebSocket service for scalability
- User authentication with JWT
- Room-based chat system
- PostgreSQL database with proper indexing
- Redis for socket clustering and caching
- Nginx reverse proxy with load balancing
- Docker containerization
- Hot reloading for development
- Production-ready nginx configuration
- Environment variable configuration

## Security Notes

- Change default passwords in production
- Use strong JWT secrets
- Configure proper CORS settings
- Enable SSL/TLS in production
- Regularly update dependencies

## Troubleshooting

### Port conflicts
If you get port conflicts, change the ports in your `.env` file:
```bash
POSTGRES_PORT=5433
BACKEND_PORT=8081
FRONTEND_PORT=3001
SOCKET_PORT=8082
NGINX_PORT=8080
```

### Database connection issues
Make sure the PostgreSQL container is running:
```bash
docker-compose logs postgres
```

### Frontend not loading
Check if the backend is running and accessible:
```bash
curl http://localhost:8080/health
```

### Permission issues
If you encounter permission issues with Docker volumes, you may need to adjust file permissions or run Docker with appropriate user permissions.
