FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY build/community-rpc /app/community-rpc
COPY app/community/rpc/etc/community-rpc.yaml /app/etc/community-rpc.yaml
EXPOSE 9003
ENTRYPOINT ["/app/community-rpc", "-f", "/app/etc/community-rpc.yaml"]
