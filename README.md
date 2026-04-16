# URL Shortener

Full-stack URL shortener with:

- `backend`: Go + Gin API (`POST /api/shorten`, `GET /api/urls`, `GET /:code`)
- `frontend`: React + Vite user interface

## Quick Start

### 1) Run backend

```bash
cd backend
go run ./cmd/server
```

Backend defaults:

- Host: `localhost`
- Port: `8080`

Optional backend environment variables:

- `HOST` (default: `localhost`)
- `PORT` (default: `8080`)
- `BASE_URL` (override short URL base in API responses)
- `CORS_ALLOW_ORIGINS` (comma-separated, default: `*`)

### 2) Run frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend dev server runs on Vite defaults and proxies `/api/*` to `http://localhost:8080`.

Optional frontend environment variable:

- `VITE_API_BASE_URL` (full backend host URL)

## Verification

Backend tests:

```bash
cd backend
GOCACHE=/tmp/go-build go test ./...
```

Frontend build:

```bash
cd frontend
npm run build
```
