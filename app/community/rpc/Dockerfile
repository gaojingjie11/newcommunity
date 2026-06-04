FROM golang:1.24-alpine AS builder

WORKDIR /build

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -p 2 -ldflags="-s -w" -o /app/community-rpc ./app/community/rpc/community.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/community-rpc /app/community-rpc
COPY app/community/rpc/etc/community-rpc.yaml /app/etc/community-rpc.yaml
EXPOSE 9003
ENTRYPOINT ["/app/community-rpc", "-f", "/app/etc/community-rpc.yaml"]
