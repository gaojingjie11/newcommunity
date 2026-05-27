# 构建阶段：编译Go代码
FROM golang:1.24-alpine AS builder
WORKDIR /app
# 复制依赖清单，利用Docker缓存
COPY go.mod go.sum ./

# 配置国内 GOPROXY 代理，加速依赖下载
RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN go mod download
# 复制所有后端代码
COPY . .
# 静态编译Go程序，适配alpine轻量镜像
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o smartserver ./cmd/main.go

# 运行阶段：轻量镜像运行
FROM alpine:latest
RUN apk add --no-cache tzdata
WORKDIR /app
# 复制编译产物
COPY --from=builder /app/smartserver .
# 复制生产环境配置文件
COPY --from=builder /app/config/prod.yaml ./config/prod.yaml
# 暴露后端API端口
EXPOSE 8082
# 启动后端服务
CMD ["./smartserver"]