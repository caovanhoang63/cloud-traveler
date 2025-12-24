# Cloud Traveler

Production-ready 3-tier architecture demo with Go, PostgreSQL, and AWS S3.

## Architecture

```
Web Tier          App Tier           Data Tier         Storage
┌─────────┐      ┌──────────┐      ┌───────────┐     ┌─────────┐
│ HTML/JS │ ──── │ Go + Gin │ ──── │ PostgreSQL│     │  AWS S3 │
└─────────┘      └──────────┘      └───────────┘     └─────────┘
```

## Features

- Database health monitoring endpoint
- File upload streaming to S3
- File metadata CRUD with PostgreSQL
- Modern responsive UI
- Multi-stage Docker build (< 20MB)
- GitHub Actions CI/CD to GHCR

## Quick Start

### Prerequisites

- Docker & Docker Compose
- AWS credentials (for S3)

### Run Locally

```bash
# Copy environment file
cp .env.example .env

# Edit .env with your AWS credentials
vim .env

# Start services
docker compose up
```

Access at http://localhost:8080

## Configuration

| Variable | Required | Default | Description |
|----------|----------|---------|-------------|
| DB_HOST | No | localhost | PostgreSQL host |
| DB_PORT | No | 5432 | PostgreSQL port |
| DB_USER | No | postgres | Database user |
| DB_PASSWORD | No | postgres | Database password |
| DB_NAME | No | cloud_traveler | Database name |
| S3_BUCKET_NAME | Yes | - | AWS S3 bucket |
| AWS_REGION | No | us-east-1 | AWS region |
| SERVER_PORT | No | 8080 | Server port |

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/` | Web UI |
| GET | `/api/health/db` | Database health check |
| POST | `/api/upload` | Upload file to S3 |
| GET | `/api/files` | List uploaded files |
| GET | `/api/files/:id` | Get file by ID |
| DELETE | `/api/files/:id` | Delete file record |

## Project Structure

```
├── main.go                 # Entry point
├── internal/
│   ├── config/            # Environment configuration
│   ├── handler/           # HTTP handlers
│   └── storage/           # Database & S3 clients
├── templates/             # HTML templates
├── docs/                  # Documentation (BRD, SRS)
├── Dockerfile             # Multi-stage build
├── docker-compose.yaml    # Local development
└── .github/workflows/     # CI/CD pipeline
```

## Development

```bash
# Install dependencies
go mod download

# Run locally (requires PostgreSQL)
go run .

# Verify code
go vet ./...
```

## Deployment

Push to `main` branch triggers GitHub Actions:
1. Build Docker image
2. Push to `ghcr.io/caovanhoang63/cloud-traveler:latest`

## Tech Stack

- **Backend:** Go 1.23, Gin
- **Database:** PostgreSQL 16
- **Storage:** AWS S3
- **Container:** Docker, Alpine Linux
- **CI/CD:** GitHub Actions, GHCR

## License

MIT
