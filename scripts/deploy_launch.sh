#!/bin/bash
set -e

DEPLOY_SERVICES=$(cat .deploy-services 2>/dev/null | xargs || true)
if [ -z "${DEPLOY_SERVICES}" ]; then
  echo "3. No deployable service changes detected, skipping container restart."
  exit 0
fi

echo "3. Restarting containers with new images: ${DEPLOY_SERVICES}"
docker-compose up -d --remove-orphans ${DEPLOY_SERVICES}

echo "4. Post-deploy docker cleanup..."
docker image prune -f || true
docker network prune -f || true

DISK_USAGE=$(df -P / | awk 'NR==2 {gsub("%","",$5); print $5}')
echo "4.1 Root disk usage after deploy: ${DISK_USAGE}%"
if [ "${DISK_USAGE}" -ge 92 ]; then
  echo "4.2 Disk usage still high, pruning old images and build cache..."
  docker image prune -af || true
  docker builder prune -af || true
fi

echo "5. Docker disk usage after deploy..."
docker system df || true

echo "Deployment successfully completed!"
