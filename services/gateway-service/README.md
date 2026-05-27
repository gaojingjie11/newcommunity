# gateway-service

Unified API entry for the microservice skeleton.

## Run

```bash
go run ./services/gateway-service/cmd/server -config services/gateway-service/configs/config.yaml
```

## Endpoints

- `GET /health`
- `GET /api/gateway/services`
- `ANY /api/proxy/:service/*path`

Nacos registration is best effort. If Nacos is unavailable, the service logs a warning and continues with local configuration.
