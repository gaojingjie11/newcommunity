# Docker Compose Deployment

## Start

```bash
cd deploy/docker-compose
docker compose up -d --build
```

## Check Services

```bash
curl http://127.0.0.1:8000/health
curl http://127.0.0.1:8000/api/gateway/services
curl http://127.0.0.1:8001/health
curl http://127.0.0.1:8002/health
curl http://127.0.0.1:8003/health
curl http://127.0.0.1:8004/health
curl http://127.0.0.1:8005/health
curl http://127.0.0.1:9100/health
```

Agent checks:

```bash
curl -X POST http://127.0.0.1:9100/agent/chat \
  -H 'Content-Type: application/json' \
  -d '{"message":"物业费怎么交"}'

curl -X POST http://127.0.0.1:9100/agent/repair-classify \
  -H 'Content-Type: application/json' \
  -d '{"content":"楼道灯坏了"}'

curl -X POST http://127.0.0.1:9100/agent/complaint-risk \
  -H 'Content-Type: application/json' \
  -d '{"content":"噪音扰民多次未解决"}'

curl -X POST http://127.0.0.1:9100/agent/recommend \
  -H 'Content-Type: application/json' \
  -d '{"scene":"home"}'
```

Consoles:

- Nacos: `http://127.0.0.1:8848/nacos`
- RabbitMQ: `http://127.0.0.1:15672`, `guest` / `guest`
- MinIO: `http://127.0.0.1:9001`, `minioadmin` / `minioadmin`

The legacy `frontend` is not included in this compose file yet. It can continue to run with its existing `frontend` workflow until API migration to `gateway-service` is started.
