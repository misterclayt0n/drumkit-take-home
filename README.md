# drumkit-take-home

First pass backend for Turvo integration.

## Structure

```text
.
├── cmd/api              # Go HTTP server entrypoint
├── internal/config      # config + godotenv loading
├── internal/httpapi     # HTTP handlers
├── internal/turvo       # Turvo auth + API client
└── frontend             # placeholder for future React app
```

This keeps the Go API and future React app in the same repo, while still treating them as separate processes.

## Environment

Create a `.env` file in the repo root:

```env
PORT=8080
BACKEND_URL=http://localhost
TURVO_BASE_URL=https://my-sandbox-publicapi.turvo.com/v1
TURVO_API=...
TURVO_CLIENT_NAME=...
TURVO_CLIENT_SECRET=...
TURVO_USERNAME=...
TURVO_PASSWORD=...
```

Notes:
- `TURVO_BASE_URL` should stay on the `publicapi` host for API calls.
- The Turvo web login URL can be different from the API base URL.

## Run

```bash
make run
```

Useful commands:

```bash
make build
make test
make fmt
make tidy
```

## Endpoints

- `GET /healthz`
- `GET /api/shipments`

Example:

```bash
curl "$BACKEND_URL:$PORT/api/shipments"
```
