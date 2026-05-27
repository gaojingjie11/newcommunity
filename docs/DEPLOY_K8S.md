# Kubernetes Deployment

The current Kubernetes files are a basic runnable skeleton, not a production high-availability setup.

## Apply Order

```bash
kubectl apply -f deploy/k8s/namespace.yaml
kubectl apply -f deploy/k8s/secrets/
kubectl apply -f deploy/k8s/configmaps/
kubectl apply -f deploy/k8s/infrastructure/
kubectl apply -f deploy/k8s/services/
kubectl apply -f deploy/k8s/ingress.yaml
```

## Notes

- Service images use placeholders such as `smartcommunity/gateway-service:latest`.
- Build and push images before applying service deployments in a real cluster.
- Persistence, probes, resource limits, autoscaling and secure secrets should be added in later stages.
- Nacos is standalone, MySQL is single pod, and RabbitMQ/MinIO are development-grade skeleton deployments.
