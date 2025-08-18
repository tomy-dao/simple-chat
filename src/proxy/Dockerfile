# Build stage
FROM nginx:alpine

# Install additional tools
RUN apk add --no-cache \
    curl \
    vim \
    bash \
    envsubst

# Create nginx user if not exists
RUN addgroup -g 1001 -S nginx && \
    adduser -u 1001 -S nginx -G nginx

# Set working directory
WORKDIR /etc/nginx

# Copy nginx configuration templates
COPY nginx.conf.template /etc/nginx/nginx.conf.template
COPY conf.d/ /etc/nginx/conf.d/

# Create necessary directories
RUN mkdir -p /var/log/nginx /var/cache/nginx /var/run

# Set proper permissions
RUN chown -R nginx:nginx /var/log/nginx /var/cache/nginx /var/run /etc/nginx

# Copy startup script
COPY docker-entrypoint.sh /
RUN chmod +x /docker-entrypoint.sh

# Expose ports
EXPOSE 80 443

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost/health || exit 1

# Start nginx with environment variable substitution
ENTRYPOINT ["/docker-entrypoint.sh"]
CMD ["nginx", "-g", "daemon off;"]
