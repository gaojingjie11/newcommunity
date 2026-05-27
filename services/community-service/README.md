# community-service

Notice, visitor, parking, property fee and basic community service.

## Run

```bash
go run ./services/community-service/cmd/server -config services/community-service/configs/config.yaml
```

## Endpoints

- `GET /health`
- `GET /api/community/ping`
- `GET /api/community/notices`
- `GET /api/community/notices/:id`
- `POST /api/community/notices/:id/read`
- `POST /api/community/visitors`
- `GET /api/community/visitors`
- `GET /api/community/parking-spaces/my`
- `PUT /api/community/parking-spaces/:id/plate`
- `GET /api/community/property-fees`
- `POST /api/community/property-fees/:id/pay`
- `GET /api/community/property-fees/payments`
- `GET /api/admin/community/notices`
- `POST /api/admin/community/notices`
- `DELETE /api/admin/community/notices/:id`
- `GET /api/admin/community/notices/:id/views`
- `GET /api/admin/community/visitors`
- `POST /api/admin/community/visitors/:id/audit`
- `GET /api/admin/community/parking-spaces`
- `POST /api/admin/community/parking-spaces`
- `POST /api/admin/community/parking-spaces/:id/assign`
- `GET /api/admin/community/parking-spaces/statistics`
- `GET /api/admin/community/property-fees`
- `POST /api/admin/community/property-fees`
- `GET /api/admin/community/property-fees/payments`

Admin routes use `RequirePermission` and shared RBAC tables. Property fee payment currently records idempotent bill status changes; real wallet deduction should be wired through mall-service/internal payment in the next stage.
