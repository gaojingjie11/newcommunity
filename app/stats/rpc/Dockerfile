FROM golang:1.25-alpine AS builder

WORKDIR /build

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -p 2 -ldflags="-s -w" -o /app/stats-rpc ./app/stats/rpc/stats.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/stats-rpc /app/stats-rpc
COPY app/stats/rpc/etc/stats-rpc.yaml /app/etc/stats-rpc.yaml
EXPOSE 9005
ENTRYPOINT ["/app/stats-rpc", "-f", "/app/etc/stats-rpc.yaml"]
