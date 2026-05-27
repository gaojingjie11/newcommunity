# Frontend Migration Plan

The existing `frontend` remains runnable. This stage only adds a lightweight `/agent` page and a service-card entry.

## Current State

- Existing API calls still target the legacy backend configuration.
- No broad frontend refactor has been performed.
- Agent page is a placeholder UI for future gateway/agent integration.

## Later Migration Steps

1. Add environment variables for `VITE_GATEWAY_BASE_URL`.
2. Switch one API module at a time from legacy backend to `gateway-service`.
3. Start with read-only endpoints: health, notices, product list and statistics overview.
4. Move auth only after `user-service` and gateway JWT middleware are implemented.
5. Keep fallback branches or flags during the transition so legacy demos continue to run.
