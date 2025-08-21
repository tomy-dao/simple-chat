# Simple Chat Application

A real-time chat application built with Go backend, React frontend, and WebSocket support, all containerized with Docker.

## Architecture

- **Backend**: Go with JWT authentication and MySQL database
- **Frontend**: React application
- **Database**: MySQL 8.0
- **Socket Service**: WebSocket support for real-time messaging
- **Containerization**: Docker Compose for easy deployment

## Features

- User registration and authentication
- JWT-based security
- Real-time messaging with WebSockets
- MySQL database integration
- Docker containerization

## Prerequisites

- Docker and Docker Compose
- Git

## Quick Start

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd simple-chat
   ```

2. **Set up environment variables**
   ```bash
   cp env.example .env
   ```
   Edit `.env` file with your desired configuration values.

3. **Start the application**
   ```bash
   docker-compose up -d
   ```

4. **Access the application**
   - Frontend: http://localhost (or your configured NGINX_PORT)
   - Backend API: http://localhost/api
   - WebSocket: ws://localhost/socket

## Environment Configuration

The application uses environment variables for configuration. Copy `env.example` to `.env` and modify the values:

### Database Configuration
- `MYSQL_DATABASE`: Database name
- `MYSQL_USER`: Database user
- `MYSQL_PASSWORD`: Database password
- `MYSQL_ROOT_PASSWORD`: Root password for MySQL

### Server Configuration
- `NGINX_PORT`: Nginx HTTP port (default: 80)
- `NGINX_SSL_PORT`: Nginx HTTPS port (default: 443)
- `BACKEND_PORT`: Backend service port (default: 8080)
- `FRONTEND_PORT`: Frontend service port (default: 3000)
- `SOCKET_PORT`: WebSocket service port (default: 8081)

### Security Configuration
- `JWT_SECRET`: Secret key for JWT tokens

## Development

### Backend Development
The Go backend is located in `src/backend/` and includes:
- JWT authentication service
- User management
- Database models and migrations

### Frontend Development
The React frontend is located in `src/frontend/` and includes:
- User interface components
- WebSocket integration
- API client

### Database
MySQL database with initialization scripts in `init.sql`.

## Docker Services

- **mysql**: MySQL 8.0 database
- **backend**: Go backend service
- **frontend**: React frontend service
- **socket**: WebSocket service for real-time messaging

