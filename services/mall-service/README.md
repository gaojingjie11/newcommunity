# mall-service

Product, category, cart, order, payment and wallet service skeleton.

## Run

```bash
go run ./services/mall-service/cmd/server -config services/mall-service/configs/config.yaml
```

## Endpoints

- `GET /health`
- `GET /api/mall/ping`
- `GET /api/mall/products`

Future migration source: legacy product, cart, order, comment, favorite, finance and transaction modules.
