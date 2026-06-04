#!/bin/bash
set -e

echo "=== 1. Starting Host Cross-Compilation for Go Services ==="
mkdir -p build

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64

echo "Building user-rpc..."
go build -ldflags="-s -w" -o ./build/user-rpc ./app/user/rpc/user.go

echo "Building mall-rpc..."
go build -ldflags="-s -w" -o ./build/mall-rpc ./app/mall/rpc/mall.go

echo "Building community-rpc..."
go build -ldflags="-s -w" -o ./build/community-rpc ./app/community/rpc/community.go

echo "Building workorder-rpc..."
go build -ldflags="-s -w" -o ./build/workorder-rpc ./app/workorder/rpc/workorder.go

echo "Building stats-rpc..."
go build -ldflags="-s -w" -o ./build/stats-rpc ./app/stats/rpc/stats.go

echo "Building gateway-api..."
go build -ldflags="-s -w" -o ./build/gateway-api ./app/gateway/api

echo "=== 2. Starting Frontend Build ==="
cd frontend
echo "Installing frontend npm packages..."
npm install --registry=https://registry.npmmirror.com
echo "Building frontend production assets..."
npm run build
cd ..

echo "=== Build Complete! You can now run: docker compose up --build -d ==="
