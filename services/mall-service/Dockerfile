FROM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY pkg ./pkg
COPY services/mall-service ./services/mall-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/mall-service ./services/mall-service/cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/mall-service /app/mall-service
COPY services/mall-service/configs/config.yaml /app/configs/config.yaml
EXPOSE 8002
ENTRYPOINT ["/app/mall-service", "-config", "/app/configs/config.yaml"]
