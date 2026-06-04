FROM golang:1.24-alpine AS builder

WORKDIR /build

ENV GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -p 2 -ldflags="-s -w" -o /app/workorder-rpc ./app/workorder/rpc/workorder.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/workorder-rpc /app/workorder-rpc
COPY app/workorder/rpc/etc/workorder-rpc.yaml /app/etc/workorder-rpc.yaml
EXPOSE 9004
ENTRYPOINT ["/app/workorder-rpc", "-f", "/app/etc/workorder-rpc.yaml"]
