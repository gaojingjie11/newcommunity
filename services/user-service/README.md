# user-service

User, login, registration, JWT, role and permission service skeleton.

## Run

```bash
go run ./services/user-service/cmd/server -config services/user-service/configs/config.yaml
```

## Endpoints

- `GET /health`
- `GET /api/users/ping`
- `POST /api/users/login`

Future migration source: legacy `UserHandler`, `UserService`, `SysUser`, role and JWT middleware.
