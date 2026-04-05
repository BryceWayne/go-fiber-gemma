# Go Fiber Gemma HTMX

A sleek, local-network API & Chat interface for serving the Gemma 4 model, featuring a hyper-modern HTMX frontend built on Go Fiber.

## Project Structure

This project follows standard Go package layout guidelines:

- `cmd/api/main.go` - The entry point to start the server.
- `internal/` - Private application and library code.
  - `handlers/` - HTTP route handlers (API & HTMX views).
  - `models/` - Data structures and payloads.
  - `service/` - Business logic and Ollama integration.
- `web/` - Frontend assets.
  - `templates/` - HTML templates used by Fiber's template engine.
  - `static/` - Static files like CSS and images.

## Features

- **Blazing Fast API**: Built with Go Fiber v2.
- **Dynamic Frontend**: Modern single-page feel without writing JavaScript, powered by HTMX.
- **Sleek Aesthetic**: Deep dark mode with glassmorphism and vector-inspired glowing accents.
- **Local AI**: Direct integration with Ollama (specifically running `gemma4:e2b`).
- **mDNS Support**: Accessible over local network via `gemma.local`.

## Getting Started

1. Ensure Ollama is running locally with the model available on port `11434`.
2. Run the server:

```bash
go run ./cmd/api
```

1. Navigate to `http://localhost:3000` to interact with the UI, or hit `http://gemma.local:3000/api/status` for API health checks.
