FROM golang:1.25-alpine AS builder

WORKDIR /build

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -p 2 -ldflags="-s -w" -o /app/gateway-api ./app/gateway/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/gateway-api /app/gateway-api
COPY app/gateway/api/etc/gateway-api.yaml /app/etc/gateway-api.yaml
EXPOSE 8000
ENTRYPOINT ["/app/gateway-api", "-f", "/app/etc/gateway-api.yaml"]
