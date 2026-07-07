#!/bin/bash
set -e

echo "0. Docker disk usage before cleanup..."
docker system df || true

echo "0.1 Lightweight pre-build cleanup..."
docker container prune -f || true
docker network prune -f || true

DISK_USAGE=$(df -P / | awk 'NR==2 {gsub("%","",$5); print $5}')
echo "0.2 Current root disk usage: ${DISK_USAGE}%"
if [ "${DISK_USAGE}" -ge 92 ]; then
  echo "0.3 Disk usage too high, running emergency docker cleanup..."
  docker image prune -af || true
  docker builder prune -af || true
  docker container prune -f || true
  docker network prune -f || true
fi

echo "0.4 Docker disk usage after cleanup..."
docker system df || true

# Calculate current .env MD5 hash to detect manual config updates on the server
ENV_CHANGED=0
if [ -f ".env" ]; then
  CURRENT_ENV_MD5=$(md5sum .env 2>/dev/null | awk '{print $1}' || md5 .env 2>/dev/null | awk '{print $1}' || true)
  if [ -n "${CURRENT_ENV_MD5}" ]; then
    PAST_ENV_MD5=$(cat .env.md5 2>/dev/null || true)
    if [ "${CURRENT_ENV_MD5}" != "${PAST_ENV_MD5}" ]; then
      echo "Config change detected: .env file has been modified."
      ENV_CHANGED=1
      echo "${CURRENT_ENV_MD5}" > .env.md5
    fi
  fi
fi

PREV_HEAD=$(git rev-parse HEAD@{1} 2>/dev/null || true)
NEW_HEAD=$(git rev-parse HEAD 2>/dev/null || true)

ALL_SERVICES="user-rpc mall-rpc community-rpc workorder-rpc stats-rpc agent-rpc gateway-api frontend"
ALL_GO_SERVICES="user-rpc mall-rpc community-rpc workorder-rpc stats-rpc agent-rpc gateway-api"
BUILD_SERVICES=""
BUILD_GO_ALL=0
BUILD_FRONTEND=0

add_service() {
  case " ${BUILD_SERVICES} " in
    *" $1 "*) ;;
    *) BUILD_SERVICES="${BUILD_SERVICES} $1" ;;
  esac
}

if [ -z "${PREV_HEAD}" ] || [ -z "${NEW_HEAD}" ] || [ "${PREV_HEAD}" = "${NEW_HEAD}" ]; then
  echo "1.1 Unable to determine changed files reliably, fallback to full build."
  BUILD_SERVICES="${ALL_SERVICES}"
else
  CHANGED_FILES=$(git diff --name-only "${PREV_HEAD}" "${NEW_HEAD}")
  echo "1.1 Changed files:"
  echo "${CHANGED_FILES}"

  for file in ${CHANGED_FILES}; do
    case "${file}" in
      go.mod|go.sum|docker-compose.yml|.dockerignore|common/*)
        BUILD_GO_ALL=1
        ;;
      app/user/rpc/*)
        add_service user-rpc
        ;;
      app/mall/rpc/*)
        add_service mall-rpc
        ;;
      app/community/rpc/*)
        add_service community-rpc
        ;;
      app/workorder/rpc/*)
        add_service workorder-rpc
        ;;
      app/stats/rpc/*)
        add_service stats-rpc
        ;;
      app/agent/rpc/*)
        add_service agent-rpc
        ;;
      app/gateway/api/*)
        add_service gateway-api
        ;;
      frontend/*)
        BUILD_FRONTEND=1
        ;;
    esac
  done

  if [ "${BUILD_GO_ALL}" -eq 1 ]; then
    BUILD_SERVICES="${ALL_GO_SERVICES}"
  fi
  if [ "${BUILD_FRONTEND}" -eq 1 ]; then
    add_service frontend
  fi
fi

BUILD_SERVICES=$(echo "${BUILD_SERVICES}" | xargs || true)

if [ -z "${BUILD_SERVICES}" ] && [ "${ENV_CHANGED}" -eq 0 ]; then
  echo "1.2 No code changes or config changes detected, skipping deploy step."
  echo "" > .deploy-services
  exit 0
fi

# Build only the services that actually had code changes
if [ -n "${BUILD_SERVICES}" ]; then
  echo "2. Services to build: ${BUILD_SERVICES}"
  echo "2. Starting limited parallel build..."
  COMPOSE_PARALLEL_LIMIT=2 docker-compose build --parallel ${BUILD_SERVICES}
  echo "All builds completed successfully!"
else
  echo "2. No code changes detected, skipping docker build."
fi

# Determine which services to launch/restart
if [ "${ENV_CHANGED}" -eq 1 ]; then
  echo "2.1 Config (.env) changed, all services will be restarted to apply variables: ${ALL_SERVICES}"
  echo "${ALL_SERVICES}" > .deploy-services
else
  echo "2.1 Only changed services will be restarted: ${BUILD_SERVICES}"
  echo "${BUILD_SERVICES}" > .deploy-services
fi
