FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY build/stats-rpc /app/stats-rpc
COPY app/stats/rpc/etc/stats-rpc.yaml /app/etc/stats-rpc.yaml
EXPOSE 9005
ENTRYPOINT ["/app/stats-rpc", "-f", "/app/etc/stats-rpc.yaml"]
