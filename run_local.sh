#!/bin/bash

# 用于保存所有后台进程的 PID
pids=()

# 退出清理函数
cleanup() {
    echo -e "\n=== 5. Stopping all local background services ==="
    for pid in "${pids[@]}"; do
        if kill -0 "$pid" 2>/dev/null; then
            kill "$pid"
            echo "Killed process $pid"
        fi
    done
    exit 0
}

# 捕获 Ctrl+C (SIGINT) 和 退出 (SIGTERM) 信号，一旦按下 Ctrl+C 则自动释放所有后台占用的端口
trap cleanup SIGINT SIGTERM

# Load environment variables if .env exists
if [ -f .env ]; then
    export $(grep -v '^#' .env | xargs)
fi

echo "=== 1. Stopping Docker Containers to Release Ports ==="
docker compose down

echo "=== 2. Starting Go-Zero RPC Microservices in Background ==="
echo "Starting user-rpc..."
go run app/user/rpc/user.go -f app/user/rpc/etc/user-rpc.yaml &
pids+=($!)
sleep 0.5

echo "Starting mall-rpc..."
go run app/mall/rpc/mall.go -f app/mall/rpc/etc/mall-rpc.yaml &
pids+=($!)
sleep 0.5

echo "Starting community-rpc..."
go run app/community/rpc/community.go -f app/community/rpc/etc/community-rpc.yaml &
pids+=($!)
sleep 0.5

echo "Starting workorder-rpc..."
go run app/workorder/rpc/workorder.go -f app/workorder/rpc/etc/workorder-rpc.yaml &
pids+=($!)
sleep 0.5

echo "Starting stats-rpc..."
go run app/stats/rpc/stats.go -f app/stats/rpc/etc/stats-rpc.yaml &
pids+=($!)
sleep 0.5

echo "Starting agent-rpc (Port 9006)..."
go run app/agent/rpc/agent.go -f app/agent/rpc/etc/agent.yaml &
pids+=($!)
sleep 0.5

echo "Starting gateway-api (Port 8000)..."
go run app/gateway/api/gateway.go -f app/gateway/api/etc/gateway-api.yaml &
pids+=($!)
sleep 1.5

echo "=== 4. Starting Vue Frontend Dev Server (Vite) ==="
cd frontend
npm run dev
