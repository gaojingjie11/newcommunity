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

PREV_HEAD=$(git rev-parse HEAD 2>/dev/null || true)

# Git pull is already run by the GitHub Actions runner before executing this script
# But we can run it again or just read the heads.
# Since we pulled in deploy.yml already, NEW_HEAD is the current HEAD.
# To compute changed files, we need the PREV_HEAD which was before we pulled.
# Wait! If git pull is run in deploy.yml, the local repository HEAD moves to NEW_HEAD.
# So PREV_HEAD in this script would be the SAME as NEW_HEAD if we run it after pulling!
# Ah! That is a very important point!
# If we run "git pull" in deploy.yml, then by the time this script runs, the HEAD has already changed!
# So git rev-parse HEAD will return the NEW head, and we won't know the PREV head!
# How do we solve this?
# We can read the git reflog on the server to see the previous commit!
# In git reflog, HEAD@{1} is the commit before the git pull!
# Yes! `git rev-parse HEAD@{1}` is exactly the PREV_HEAD!
# This is a brilliant and standard git trick!
# Let's verify: does HEAD@{1} point to the commit before pulling?
# Yes, because git pull updates the local head, which appends to the reflog, shifting the previous head to HEAD@{1}!

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
echo "${BUILD_SERVICES}" > .deploy-services

if [ -z "${BUILD_SERVICES}" ]; then
  echo "1.2 No deployable service changes detected, skipping build."
  exit 0
fi

echo "2. Services to build: ${BUILD_SERVICES}"
echo "2. Starting limited parallel build..."
COMPOSE_PARALLEL_LIMIT=2 docker-compose build --parallel ${BUILD_SERVICES}
echo "All builds completed successfully!"
