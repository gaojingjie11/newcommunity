# Smart Community Microservices Skeleton

This repository now keeps the original `frontend` and `smartcomunity` monolith as legacy reference code, and adds a new B-plan microservice skeleton.

## Layout

```text
pkg/                         shared Go packages
services/
  gateway-service            Go API gateway skeleton
  user-service               Go user service skeleton
  mall-service               Go mall/order/payment skeleton
  community-service          Go community service skeleton
  workorder-service          Go repair/complaint skeleton with RabbitMQ event placeholders
  statistics-service         Go statistics skeleton
  agent-service              Python FastAPI agent skeleton
deploy/docker-compose        local infrastructure and service orchestration
deploy/k8s                   Kubernetes skeleton
docs                         architecture, deployment and migration handoff docs
```

## Quick Start

```bash
cd deploy/docker-compose
docker compose up -d --build
```

Main checks:

```bash
curl http://127.0.0.1:8000/health
curl http://127.0.0.1:8000/api/gateway/services
curl http://127.0.0.1:9100/health
curl -X POST http://127.0.0.1:9100/agent/chat -H 'Content-Type: application/json' -d '{"message":"hello"}'
```

Consoles:

- Nacos: `http://127.0.0.1:8848/nacos`
- RabbitMQ: `http://127.0.0.1:15672` (`guest` / `guest`)
- MinIO Console: `http://127.0.0.1:9001` (`minioadmin` / `minioadmin`)

MinIO API is mapped to host `19000` because `agent-service` and MinIO both use container port `9000`; inside Docker network, services still use `minio:9000` and `agent-service:9000`.

See `docs/MIGRATION_PROGRESS.md` for current status and next handoff tasks.
