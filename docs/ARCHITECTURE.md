# Architecture

The current stage creates a microservice framework beside the legacy monolith. The original `smartcomunity` Go application and `frontend` remain intact as reference implementations.

```text
Browser / frontend
       |
       v
gateway-service :8000
       |
       +--> user-service :8001
       +--> mall-service :8002
       +--> community-service :8003
       +--> workorder-service :8004 -- publishes --> RabbitMQ
       +--> statistics-service :8005
       +--> agent-service :9000

Shared infrastructure:
Nacos       service registry and future config center
MySQL       shared development database, smart_community
Redis       cache, captcha, session and lightweight distributed state
RabbitMQ    domain events: repair.created, complaint.created, order.created
MinIO       images, attachments, repair photos and product images
```

Go services share packages under `pkg/` for config, response shape, middleware, auth, database, Redis, Nacos, RabbitMQ and MinIO. Nacos registration is best effort: unavailable Nacos logs a warning and does not stop service startup.

`agent-service` is a Python FastAPI service for future LLM workflows. It currently returns deterministic placeholder responses and reserves `LLM_API_KEY`, `LLM_BASE_URL`, and `LLM_MODEL`.
