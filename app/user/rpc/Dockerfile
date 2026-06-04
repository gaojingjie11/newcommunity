FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY build/user-rpc /app/user-rpc
COPY app/user/rpc/etc/user-rpc.yaml /app/etc/user-rpc.yaml
EXPOSE 9001
ENTRYPOINT ["/app/user-rpc", "-f", "/app/etc/user-rpc.yaml"]
