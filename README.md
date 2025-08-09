<h1 align="center">Unfurl</h1>

## :page_facing_up: Description

Fast and reliable URL metadata extraction service.

## :hammer_and_wrench: Development

- Go 1.24+
- Docker
- Make

## :sparkles: Features

- Smart fallback strategy (HTTP → Playwright)
- Handles static sites and JS-heavy applications
- Extracts titles, descriptions, images, and site names
- Containerized with multi-stage builds

## :rocket: Quick Start

### Using Docker

```bash
# Pull and run
docker run -p 8080:8080 mtsfy/unfurl:latest

# Or build locally
git clone https://github.com/mtsfy/unfurl.git
cd unfurl
make docker-build
make docker-run
```

## :clipboard: API Reference

### Extract URL Metadata

```http
POST /api/v1/unfurl
Content-Type: application/json
{
  "url": "https://github.com/mtsfy/linkself"
}
```

**Response:**

```json
{
  "title": "GitHub - mtsfy/linkself: Application to manage multiple links.",
  "description": "Application to manage multiple links. Contribute to mtsfy/linkself development by creating an account on GitHub.",
  "image": "https://opengraph.githubassets.com/00bc6752871fc260d420483a2c30a3050005b97d4fe33fa6f259c6898b78815a/mtsfy/linkself",
  "site": "GitHub"
}
```

### Health Check

```http
GET /api/v1/health
```

```json
{
  "status": "OK!",
  "timestamp": "2025-08-09T17:30:00Z"
}
```
