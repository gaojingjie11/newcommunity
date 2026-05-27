FROM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY pkg ./pkg
COPY services/gateway-service ./services/gateway-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/gateway-service ./services/gateway-service/cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/gateway-service /app/gateway-service
COPY services/gateway-service/configs/config.yaml /app/configs/config.yaml
EXPOSE 8000
ENTRYPOINT ["/app/gateway-service", "-config", "/app/configs/config.yaml"]
