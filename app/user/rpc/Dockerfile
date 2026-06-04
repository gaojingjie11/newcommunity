FROM golang:1.24-alpine AS builder

WORKDIR /build

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -p 2 -ldflags="-s -w" -o /app/user-rpc ./app/user/rpc/user.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/user-rpc /app/user-rpc
COPY app/user/rpc/etc/user-rpc.yaml /app/etc/user-rpc.yaml
EXPOSE 9001
ENTRYPOINT ["/app/user-rpc", "-f", "/app/etc/user-rpc.yaml"]
