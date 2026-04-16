# URL Shortener Frontend

React + TypeScript frontend for the URL shortener service.

## Features

- Create short URLs from long links.
- View recent links with click counts and creation timestamps.
- Filter and sort recent links (`newest`, `oldest`, `most clicked`).
- Copy latest generated short URL to clipboard.

## Development

```bash
npm install
npm run dev
```

The Vite dev server proxies `/api/*` requests to `http://localhost:8080` by default.

## Configuration

You can point the frontend to a different backend URL by setting:

- `VITE_API_BASE_URL`: full base URL of the backend API host (for example: `https://short.example.com`).

When unset, API requests are sent to relative paths (`/api/...`).

## Build

```bash
npm run build
npm run preview
```
