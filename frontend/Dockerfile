# Build stage
FROM node:18.19.1-alpine AS builder

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the React app
RUN npm run build

# Production stage with Nginx
FROM nginx:1.25-alpine

# Copy built files from builder stage
COPY --from=builder /app/build /usr/share/nginx/html

# Copy nginx configuration for Kubernetes
COPY nginx.k8s.conf /etc/nginx/conf.d/default.conf

# Expose port 80
EXPOSE 80

# Nginx will automatically process templates in /etc/nginx/templates/
# and substitute environment variables using envsubst
CMD ["nginx", "-g", "daemon off;"]
