# frontend

TanStack Start frontend for the Drumkit × Turvo take-home.

## Run

```bash
npm install
cp .env.example .env
npm run dev
```

The app runs on `http://localhost:3000` and expects the Go backend at:

```env
VITE_API_BASE_URL=http://localhost:6969
```

## Features

- lists loads from `GET /v1/loads`
- creates loads through `POST /v1/integrations/webhooks/loads`
- styled with Tailwind + shadcn/ui
- uses TanStack Router, Query, and Table

## Build

```bash
npm run build
```
