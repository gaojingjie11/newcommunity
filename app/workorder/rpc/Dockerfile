FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY build/workorder-rpc /app/workorder-rpc
COPY app/workorder/rpc/etc/workorder-rpc.yaml /app/etc/workorder-rpc.yaml
EXPOSE 9004
ENTRYPOINT ["/app/workorder-rpc", "-f", "/app/etc/workorder-rpc.yaml"]
