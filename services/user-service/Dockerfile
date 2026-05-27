FROM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY pkg ./pkg
COPY services/user-service ./services/user-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/user-service ./services/user-service/cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/user-service /app/user-service
COPY services/user-service/configs/config.yaml /app/configs/config.yaml
EXPOSE 8001
ENTRYPOINT ["/app/user-service", "-config", "/app/configs/config.yaml"]
