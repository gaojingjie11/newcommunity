# agent-service

FastAPI skeleton for future LLM-powered agents.

## Run

```bash
uvicorn app.main:app --host 0.0.0.0 --port 9000
```

## Endpoints

- `GET /health`
- `GET /agent/health`
- `POST /agent/chat`
- `POST /agent/repair-classify`
- `POST /agent/complaint-risk`
- `POST /agent/recommend`

The current implementation returns deterministic placeholder JSON. Environment variables `LLM_API_KEY`, `LLM_BASE_URL`, and `LLM_MODEL` are reserved for the next stage.
