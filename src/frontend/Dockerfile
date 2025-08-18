# Build stage
FROM node:22.17-alpine AS builder

# Set working directory
WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN yarn

# Copy source code
COPY . .

ARG VITE_API_BASE_API_URL
ARG VITE_SOCKET_BASE_URL
ENV VITE_API_BASE_API_URL=$VITE_API_BASE_API_URL
ENV VITE_SOCKET_BASE_URL=$VITE_SOCKET_BASE_URL

# Build the application
RUN yarn build

# Production stage
FROM nginx:alpine

# Copy built app from builder stage
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 3000

# Start nginx
CMD ["nginx", "-g", "daemon off;"]
