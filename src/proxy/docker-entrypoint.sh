#!/bin/bash
set -e

# Function to substitute environment variables in nginx config
substitute_env_vars() {
    echo "Substituting environment variables in nginx configuration..."
    
    # Substitute environment variables in nginx.conf.template
    envsubst '${BACKEND_PORT} ${FRONTEND_PORT} ${SOCKET_PORT}' < /etc/nginx/nginx.conf.template > /etc/nginx/nginx.conf
    
    # Substitute environment variables in conf.d files
    for file in /etc/nginx/conf.d/*.conf; do
        if [ -f "$file" ]; then
            envsubst '${BACKEND_PORT} ${FRONTEND_PORT} ${SOCKET_PORT}' < "$file" > "$file.tmp"
            mv "$file.tmp" "$file"
        fi
    done
    
    echo "Environment variable substitution completed."
}

# Function to test nginx configuration
test_nginx_config() {
    echo "Testing nginx configuration..."
    if nginx -t; then
        echo "Nginx configuration is valid."
        return 0
    else
        echo "Nginx configuration is invalid!"
        return 1
    fi
}

# Function to handle shutdown
cleanup() {
    echo "Shutting down nginx..."
    nginx -s quit
    exit 0
}

# Set up signal handlers
trap cleanup SIGTERM SIGINT

# Main execution
echo "Starting nginx proxy service..."

# Substitute environment variables
substitute_env_vars

# Test configuration
if ! test_nginx_config; then
    echo "Failed to start nginx due to invalid configuration."
    exit 1
fi

# Start nginx
echo "Starting nginx..."
exec "$@"
