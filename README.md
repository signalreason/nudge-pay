# NudgePay

Automated invoice reminder SaaS for freelancers and small agencies.

## Purpose
- Help freelancers and small agencies automate invoice reminders.

## Goals
- Schedule multi-stage reminders relative to invoice due dates.
- Provide a backend API, worker, and frontend for managing clients and invoices.
- Run locally on SQLite with a clear path to Postgres at scale.

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

See `docs/PRODUCT.md` for product analysis, roadmap, and architecture. Repo map: `docs/REPO_MAP.md`.
