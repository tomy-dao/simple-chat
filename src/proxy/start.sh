#!/bin/bash

# Function to reload nginx
reload_nginx() {
    echo "Reloading nginx configuration..."
    nginx -s reload
}

# Function to check if nginx is running
check_nginx() {
    if ! pgrep -x "nginx" > /dev/null; then
        echo "Nginx is not running. Starting nginx..."
        nginx -g "daemon off;" &
    fi
}

# Start nginx in background
echo "Starting nginx..."
nginx -g "daemon off;" &

# Store nginx PID
NGINX_PID=$!

# Function to handle shutdown
cleanup() {
    echo "Shutting down nginx..."
    kill $NGINX_PID
    exit 0
}

# Set up signal handlers
trap cleanup SIGTERM SIGINT

# Monitor configuration files for changes
echo "Monitoring nginx configuration files for changes..."
while true; do
    # Check if nginx is still running
    if ! kill -0 $NGINX_PID 2>/dev/null; then
        echo "Nginx process died. Exiting..."
        exit 1
    fi
    
    # Test nginx configuration
    if nginx -t > /dev/null 2>&1; then
        # Configuration is valid, no need to reload
        sleep 5
    else
        echo "Invalid nginx configuration detected. Skipping reload."
        sleep 5
    fi
done
