# Repo Map: nudge-pay

## Purpose and scope
Invoice reminder SaaS with Go backend, Next.js frontend, and supporting infra/worker components.

## Quickstart commands
- Backend: `cd backend && NUDGEPAY_JWT_SECRET=dev-secret NUDGEPAY_DB=./nudgepay.db go run ./cmd/server`
- Frontend: `cd frontend && NEXT_PUBLIC_API_URL=http://localhost:8080 npm run dev`
- Docker: `docker compose up --build`
- Tests: `cd backend && go test ./...` and `cd frontend && npm test`

## Top-level map
- `backend/` - Go API server and core business logic.
  - `backend/cmd/server/` - API entry point.
  - `backend/internal/` - domain logic, reminders, and persistence.
  - `backend/openapi.yaml` - API contract.
- `frontend/` - Next.js app.
  - `frontend/app/` - routes/pages.
  - `frontend/components/` - UI components.
  - `frontend/tests/` - frontend tests.
- `docs/` - product and architecture notes.
- `infra/` - deployment manifests (k8s).
- `docker-compose.yml` - local multi-service dev.
- `README.md` - quickstart and tests.

## Key entry points
- `backend/cmd/server/` - Go API server.
- `backend/openapi.yaml` - API source of truth.
- `frontend/app/` - Next.js routes.
- `docker-compose.yml` - local orchestration.

## Core flows and data movement
- Frontend calls Go API -> persistence in SQLite/Postgres.
- Reminder schedules stored as offsets from due date (see docs).
- Worker/outbox processes reminder scheduling (implementation in backend).

## External integrations
- Database: SQLite (local) or Postgres (scale).
- Email/SMS provider integration not documented here.

## Configuration and deployment
- Env vars: `NUDGEPAY_JWT_SECRET`, `NUDGEPAY_DB`, `NEXT_PUBLIC_API_URL`.
- Dockerfiles under `backend/` and `frontend/`.
- Kubernetes manifests under `infra/k8s`.

## Common workflows (build/test/release)
- `docker compose up --build`
- `cd backend && go test ./...`
- `cd frontend && npm test`

## Read-next list
- `README.md`
- `docs/PRODUCT.md`
- `backend/openapi.yaml`
- `backend/cmd/server/`
- `frontend/app/`
- `infra/k8s/`

## Unknowns and follow-ups
- Worker entry point and deployment process are not spelled out in README.
