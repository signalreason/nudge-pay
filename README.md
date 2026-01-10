# NudgePay

Automated invoice reminder SaaS for freelancers and small agencies.

## Quick start

```bash
# Backend
cd backend
NUDGEPAY_JWT_SECRET=dev-secret NUDGEPAY_DB=./nudgepay.db go run ./cmd/server

# Frontend
cd frontend
NEXT_PUBLIC_API_URL=http://localhost:8080 npm run dev
```

## Docker

```bash
docker compose up --build
```

## Tests

```bash
cd backend && go test ./...
cd frontend && npm test
```

See `docs/PRODUCT.md` for product analysis, roadmap, and architecture.
