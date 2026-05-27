# Infrastructure

All defaults are for local course/development environments only. Change passwords and persistence settings before production use.

## Nacos

- Purpose: service registration and future config center
- Compose URL: `http://127.0.0.1:8848/nacos`
- Mode: standalone, auth disabled
- Note: current services use best-effort HTTP registration. If Nacos is down, services keep running.

## MySQL

- Image: `mysql:8.0`
- Database: `smart_community`
- Root password: `root123456`
- Compose host port: `3306`
- Init SQL: `deploy/docker-compose/mysql/init`

## Redis

- Image: `redis:7-alpine`
- Compose host port: `6379`
- Purpose: cache, captcha, session and simple distributed state

## RabbitMQ

- Image: `rabbitmq:3-management`
- AMQP: `127.0.0.1:5672`
- Management UI: `http://127.0.0.1:15672`
- Default login: `guest` / `guest`
- Reserved queues/events: `repair.created`, `complaint.created`, `order.created`

## MinIO

- Console: `http://127.0.0.1:9001`
- Login: `minioadmin` / `minioadmin`
- Internal endpoint: `minio:9000`
- Host API endpoint: `http://127.0.0.1:19000`
- Purpose: images, attachments, repair photos and product images

MinIO API maps to host `19000` because the agent service also uses container port `9000`.
