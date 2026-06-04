FROM golang:1.24-alpine AS builder

WORKDIR /build

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -p 2 -ldflags="-s -w" -o /app/mall-rpc ./app/mall/rpc/mall.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/mall-rpc /app/mall-rpc
COPY app/mall/rpc/etc/mall-rpc.yaml /app/etc/mall-rpc.yaml
EXPOSE 9002
ENTRYPOINT ["/app/mall-rpc", "-f", "/app/etc/mall-rpc.yaml"]
