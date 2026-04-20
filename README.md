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
LOAD_STORE_PATH=.data/drumkit-loads.json
```

Notes:
- `TURVO_BASE_URL` should stay on the `publicapi` host for API calls.
- The Turvo web login URL can be different from the API base URL.
- `LOAD_STORE_PATH` preserves the exact Drumkit load payload for created loads so list responses can mirror the Drumkit schema even when Turvo lacks native fields.

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
- `GET /v1/loads`
- `POST /v1/integrations/webhooks/loads`

List example:

```bash
curl "$BACKEND_URL:$PORT/v1/loads"
```

Optional query params:

- `status`
- `customerId`
- `pickupDateSearchFrom` (YYYY-MM-DD)
- `pickupDateSearchTo` (YYYY-MM-DD)
- `page`
- `limit`

Create example:

```bash
curl -X POST "$BACKEND_URL:$PORT/v1/integrations/webhooks/loads" \
  -H 'Content-Type: application/json' \
  -d '{
    "status": "Tendered",
    "customer": {
      "externalTMSId": "834045",
      "name": "37th St Bakery"
    },
    "pickup": {
      "externalTMSId": "525513",
      "name": "1611 CGT- Rockingham (North Carolina)",
      "city": "ROCKINGHAM",
      "state": "NC",
      "country": "US",
      "apptTime": "2026-05-01T14:00:00Z",
      "timezone": "America/New_York"
    },
    "consignee": {
      "externalTMSId": "525541",
      "name": "AAA TRANS WORLD EXPRESS, INC",
      "city": "JAMAICA",
      "state": "NY",
      "country": "US",
      "apptTime": "2026-05-02T18:00:00Z",
      "timezone": "America/New_York"
    },
    "poNums": "TEST-PO-1"
  }'
```
