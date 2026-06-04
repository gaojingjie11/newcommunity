FROM golang:1.24-alpine AS builder

WORKDIR /build

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -p 2 -ldflags="-s -w" -o /app/gateway-api ./app/gateway/api/gateway.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/gateway-api /app/gateway-api
COPY app/gateway/api/etc/gateway-api.yaml /app/etc/gateway-api.yaml
EXPOSE 8000
ENTRYPOINT ["/app/gateway-api", "-f", "/app/etc/gateway-api.yaml"]
