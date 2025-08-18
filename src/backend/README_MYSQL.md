# MySQL Setup Instructions

This project has been converted from MongoDB to MySQL. Follow these steps to set up and run the application.

## Prerequisites

1. MySQL Server (version 5.7 or higher) OR Docker
2. Go 1.23.4 or higher

## Quick Start with Docker (Recommended)

1. **Start MySQL with Docker Compose**
   ```bash
   docker-compose up -d
   ```
   
   This will start:
   - MySQL 8.0 on port 3306
   - phpMyAdmin on port 8081 (http://localhost:8081)
   
   Default credentials:
   - Root password: `password`
   - Database: `local_db`
   - User: `local_user` / Password: `local_password`

2. **Set Environment Variables**
   ```bash
   export DB_HOST=localhost
   export DB_PORT=3306
   export DB_USER=root
   export DB_PASSWORD=password
   export DB_NAME=local_db
   ```

3. **Run the Application**
   ```bash
   go run main.go
   ```

## Manual Database Setup

1. **Install MySQL Server**
   ```bash
   # On macOS with Homebrew
   brew install mysql
   
   # On Ubuntu/Debian
   sudo apt-get install mysql-server
   
   # On CentOS/RHEL
   sudo yum install mysql-server
   ```

2. **Start MySQL Service**
   ```bash
   # On macOS
   brew services start mysql
   
   # On Linux
   sudo systemctl start mysql
   ```

3. **Create Database**
   ```bash
   mysql -u root -p
   ```
   
   In MySQL console:
   ```sql
   CREATE DATABASE local_db;
   USE local_db;
   ```

4. **Run Migration** (Optional - GORM will auto-migrate)
   ```bash
   mysql -u root -p local_db < migrations/001_create_locals_table.sql
   ```

## Environment Configuration

Create a `.env` file in the project root with the following variables:

```env
# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your_password
DB_NAME=local_db

# Server Configuration
HTTP_PORT=80
GRPC_PORT=8080
HOST=0.0.0.0
```

## Running the Application

1. **Install Dependencies**
   ```bash
   go mod tidy
   ```

2. **Run the Application**
   ```bash
   go run main.go
   ```

## Key Changes Made

1. **Model Changes**: 
   - Replaced MongoDB BSON tags with GORM tags
   - Changed ID field from `primitive.ObjectID` to `uint` for auto-increment
   - Added timestamps and soft delete support
   - Changed `Role` from `[]string` to `string` (JSON) for MySQL compatibility

2. **Database Driver**: 
   - Replaced `go.mongodb.org/mongo-driver` with `gorm.io/driver/mysql`
   - Added GORM ORM for database operations

3. **Configuration**: 
   - Updated default database settings for MySQL
   - Changed default port from 5432 (PostgreSQL) to 3306 (MySQL)

4. **Repository Layer**: 
   - Created MySQL-specific repository implementation
   - Added proper error handling and logging
   - Implemented CRUD operations using GORM

## Database Schema

The `locals` table includes the following fields:
- `id`: Auto-incrementing primary key
- `local_id`: Unique identifier
- `local_name`: Name field
- `email`: Email address (unique)
- `password`: Password field
- `phone_number`: Phone number
- `role`: JSON field for roles
- `status`: Status field
- `active`: Boolean flag
- `created_at`: Creation timestamp
- `updated_at`: Last update timestamp
- `deleted_at`: Soft delete timestamp

## Docker Commands

```bash
# Start services
docker-compose up -d

# Stop services
docker-compose down

# View logs
docker-compose logs mysql

# Access MySQL directly
docker exec -it local-mysql mysql -u root -p

# Access phpMyAdmin
# Open http://localhost:8081 in your browser
```

## Troubleshooting

1. **Connection Issues**: Ensure MySQL is running and accessible
2. **Permission Issues**: Make sure the database user has proper permissions
3. **Migration Issues**: GORM will auto-migrate the schema, but you can run the SQL migration manually if needed
4. **Docker Issues**: 
   - If port 3306 is already in use, change the port mapping in docker-compose.yml
   - If containers fail to start, check logs with `docker-compose logs`
