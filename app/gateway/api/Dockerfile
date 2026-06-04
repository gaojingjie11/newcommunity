FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY build/gateway-api /app/gateway-api
COPY app/gateway/api/etc/gateway-api.yaml /app/etc/gateway-api.yaml
EXPOSE 8000
ENTRYPOINT ["/app/gateway-api", "-f", "/app/etc/gateway-api.yaml"]
