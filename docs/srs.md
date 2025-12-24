# Software Requirements Specification (SRS)

## Cloud Traveler - 3-Tier Architecture Demo

**Version:** 1.0
**Date:** 2025-12-25
**Status:** Approved

---

## 1. Introduction

### 1.1 Purpose
This document specifies software requirements for Cloud Traveler, a Go-based web application demonstrating 3-tier architecture with PostgreSQL and AWS S3 integration.

### 1.2 Scope
Cloud Traveler provides database health monitoring, file uploads to S3, and file metadata management through a modern web interface.

### 1.3 Definitions

| Term | Definition |
|------|------------|
| 3-Tier | Web (UI) + App (Logic) + Data (Storage) architecture |
| S3 | AWS Simple Storage Service |
| GHCR | GitHub Container Registry |
| Gin | Go web framework |

## 2. System Architecture

```
┌─────────────────────────────────────────────────────────┐
│                      Web Tier                           │
│                   (HTML/CSS/JS)                         │
└────────────────────────┬────────────────────────────────┘
                         │ HTTP
┌────────────────────────▼────────────────────────────────┐
│                      App Tier                           │
│                  (Go + Gin Framework)                   │
│  ┌──────────┐  ┌──────────────┐  ┌─────────────────┐   │
│  │ Handlers │  │  Repository  │  │   S3 Client     │   │
│  └──────────┘  └──────────────┘  └─────────────────┘   │
└────────┬───────────────┬────────────────┬───────────────┘
         │               │                │
┌────────▼───────┐ ┌─────▼─────┐ ┌────────▼────────┐
│   PostgreSQL   │ │  AWS S3   │ │  Config (ENV)   │
│   (Data Tier)  │ │ (Storage) │ │                 │
└────────────────┘ └───────────┘ └─────────────────┘
```

## 3. Functional Requirements

### 3.1 Database Health Check (FR-01)

| Attribute | Specification |
|-----------|---------------|
| Endpoint | `GET /api/health/db` |
| Function | Ping PostgreSQL, return status |
| Response | `{"status": "success/fail", "message": "..."}` |

### 3.2 File Upload (FR-02)

| Attribute | Specification |
|-----------|---------------|
| Endpoint | `POST /api/upload` |
| Input | `multipart/form-data` with `file` field |
| Function | Stream file to S3, save metadata to DB |
| Response | `{"status": "success", "data": {...}}` |

### 3.3 List Files (FR-03)

| Attribute | Specification |
|-----------|---------------|
| Endpoint | `GET /api/files` |
| Query Params | `limit` (default: 20), `offset` (default: 0) |
| Response | `{"status": "success", "data": {"files": [...], "total": n}}` |

### 3.4 Get File (FR-04)

| Attribute | Specification |
|-----------|---------------|
| Endpoint | `GET /api/files/:id` |
| Response | `{"status": "success", "data": {...}}` |

### 3.5 Delete File (FR-05)

| Attribute | Specification |
|-----------|---------------|
| Endpoint | `DELETE /api/files/:id` |
| Function | Remove file record from database |
| Response | `{"status": "success", "s3_key": "..."}` |

### 3.6 Web Interface (FR-06)

| Component | Specification |
|-----------|---------------|
| Architecture Diagram | Visual 4-box layout |
| DB Health Button | AJAX call to `/api/health/db` |
| Upload Form | Drag-drop + click file selection |
| Files Table | List with ID, filename, type, size, date, actions |

## 4. Non-Functional Requirements

### 4.1 Performance (NFR-01)

| Metric | Requirement |
|--------|-------------|
| Health check response | < 100ms |
| File upload | Streaming (no disk buffer) |
| Max upload size | 32MB |
| Connection pool | 25 max, 5 idle |

### 4.2 Deployment (NFR-02)

| Metric | Requirement |
|--------|-------------|
| Docker image size | < 20MB |
| Base image | alpine:latest |
| Build | Multi-stage (golang:alpine → alpine) |

### 4.3 Configuration (NFR-03)

| Variable | Required | Default |
|----------|----------|---------|
| DB_HOST | No | localhost |
| DB_PORT | No | 5432 |
| DB_USER | No | postgres |
| DB_PASSWORD | No | postgres |
| DB_NAME | No | cloud_traveler |
| S3_BUCKET_NAME | Yes | - |
| AWS_REGION | No | us-east-1 |
| SERVER_PORT | No | 8080 |

### 4.4 Security (NFR-04)

- No hardcoded credentials
- Environment-based configuration
- SQL injection prevention via parameterized queries
- Input validation on file uploads

## 5. Database Schema

### 5.1 uploaded_files

| Column | Type | Constraints |
|--------|------|-------------|
| id | SERIAL | PRIMARY KEY |
| filename | VARCHAR(255) | NOT NULL |
| s3_key | VARCHAR(255) | NOT NULL, UNIQUE |
| content_type | VARCHAR(100) | - |
| size_bytes | BIGINT | - |
| created_at | TIMESTAMP | DEFAULT CURRENT_TIMESTAMP |

**Index:** `idx_uploaded_files_s3_key` on `s3_key`

## 6. API Response Format

### 6.1 Success Response
```json
{
  "status": "success",
  "message": "Optional message",
  "data": { }
}
```

### 6.2 Error Response
```json
{
  "status": "fail",
  "message": "Error description"
}
```

## 7. Project Structure

```
cloud-traveler/
├── main.go                     # Entry point
├── internal/
│   ├── config/
│   │   └── config.go          # Environment config
│   ├── handler/
│   │   ├── health-handler.go  # DB health endpoint
│   │   ├── upload-handler.go  # S3 upload endpoint
│   │   └── file-handler.go    # File CRUD endpoints
│   └── storage/
│       ├── database.go        # PostgreSQL connection
│       ├── file-repository.go # File CRUD operations
│       └── s3-client.go       # AWS S3 client
├── templates/
│   └── index.html             # Web UI
├── Dockerfile                  # Multi-stage build
├── docker-compose.yaml        # Local development
└── .github/workflows/
    └── deploy.yml             # CI/CD pipeline
```

## 8. Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| github.com/gin-gonic/gin | v1.11.0 | Web framework |
| github.com/lib/pq | v1.10.9 | PostgreSQL driver |
| github.com/aws/aws-sdk-go-v2 | v1.41.0 | AWS SDK |
| github.com/google/uuid | v1.6.0 | UUID generation |

## 9. CI/CD Pipeline

### 9.1 Trigger
Push to `main` branch

### 9.2 Steps
1. Checkout repository
2. Setup Docker Buildx
3. Login to GHCR
4. Build Docker image
5. Push to `ghcr.io/{repository}:latest`

### 9.3 Permissions
- `contents: read`
- `packages: write`

## 10. Testing Checklist

| Test Case | Expected Result |
|-----------|-----------------|
| DB health (connected) | Status 200, success |
| DB health (disconnected) | Status 503, fail |
| Upload valid file | Status 200, file in S3 + DB |
| Upload no file | Status 400, error message |
| List files (empty) | Status 200, empty array |
| List files (with data) | Status 200, file array |
| Delete existing file | Status 200, s3_key returned |
| Delete non-existent | Status 404, error |

---

**Approval:** Document approved for development.
