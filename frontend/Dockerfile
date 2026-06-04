# Build stage
FROM node:20-alpine AS builder
WORKDIR /build

# Use Alibaba NPM mirror for faster installation in China
RUN npm config set registry https://registry.npmmirror.com

# Copy dependency definitions
COPY package*.json ./
RUN npm install

# Copy source code and build
COPY . .
RUN npm run build

# Production stage
FROM nginx:alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=builder /build/dist /usr/share/nginx/html
EXPOSE 81
CMD ["nginx", "-g", "daemon off;"]
