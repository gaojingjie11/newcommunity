FROM golang:1.24-alpine AS build

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY pkg ./pkg
COPY services/community-service ./services/community-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/community-service ./services/community-service/cmd/server

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=build /out/community-service /app/community-service
COPY services/community-service/configs/config.yaml /app/configs/config.yaml
EXPOSE 8003
ENTRYPOINT ["/app/community-service", "-config", "/app/configs/config.yaml"]
