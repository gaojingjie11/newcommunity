FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY build/mall-rpc /app/mall-rpc
COPY app/mall/rpc/etc/mall-rpc.yaml /app/etc/mall-rpc.yaml
EXPOSE 9002
ENTRYPOINT ["/app/mall-rpc", "-f", "/app/etc/mall-rpc.yaml"]
