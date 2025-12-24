# Business Requirements Document (BRD)

## Cloud Traveler - 3-Tier Architecture Demo

**Version:** 1.0
**Date:** 2025-12-25
**Status:** Approved

---

## 1. Executive Summary

Cloud Traveler is a demonstration web application showcasing production-ready 3-tier architecture patterns. The application serves as a reference implementation for cloud-native development practices, featuring PostgreSQL database integration and AWS S3 storage capabilities.

## 2. Target Market

| Attribute | Value |
|-----------|-------|
| Region | Vietnam |
| Expected Active Users | ~100 |
| User Type | Developers, learners, small teams |
| Language | English (UI) |

### 2.1 User Profile
- Vietnamese developers learning cloud architecture
- Small development teams needing file upload reference
- Students studying 3-tier application design

### 2.2 Infrastructure Considerations
- AWS Region: `ap-southeast-1` (Singapore) recommended for Vietnam users
- Low-traffic design suitable for ~100 concurrent users
- Single instance deployment sufficient

## 3. Business Objectives

| Objective | Description |
|-----------|-------------|
| Educational | Demonstrate 3-tier architecture (Web, App, Data/Storage) |
| Cloud-Native | Showcase AWS S3 integration patterns |
| DevOps Ready | Provide containerized deployment with CI/CD |
| Production Patterns | Implement real-world coding practices |

## 3. Scope

### 3.1 In Scope

- Single-page web interface with modern UI
- PostgreSQL database connectivity with health monitoring
- AWS S3 file upload functionality
- File metadata management (CRUD operations)
- Docker containerization
- GitHub Actions CI/CD pipeline
- Environment-based configuration

### 3.2 Out of Scope

- User authentication/authorization
- File download from S3
- Multi-tenant support
- Production-grade logging/monitoring
- Load balancing configuration

## 4. Stakeholders

| Role | Responsibility |
|------|----------------|
| Developers | Reference implementation for cloud architecture |
| DevOps Engineers | Deployment patterns and CI/CD reference |
| Solution Architects | 3-tier design patterns |
| Students/Learners | Educational resource |

## 5. Business Requirements

### BR-01: Database Connectivity
**Priority:** High
Application must connect to PostgreSQL and provide health status endpoint.

### BR-02: Cloud Storage Integration
**Priority:** High
Application must upload files directly to AWS S3 buckets.

### BR-03: File Management
**Priority:** Medium
Application must track uploaded files with metadata in database.

### BR-04: Container Deployment
**Priority:** High
Application must be deployable via Docker with image size under 20MB.

### BR-05: CI/CD Pipeline
**Priority:** Medium
Application must have automated build and push to container registry.

## 6. Success Criteria

| Metric | Target |
|--------|--------|
| Docker image size | < 20MB |
| Database health check | < 100ms response |
| File upload | Stream directly to S3 |
| Build pipeline | Auto-trigger on master branch |
| Concurrent users | Support ~100 active users |
| Availability | 99% uptime (single instance) |

## 7. Constraints

- Must use Go programming language
- Must use Gin web framework
- Must use AWS SDK v2 for S3
- Must use lib/pq for PostgreSQL
- Configuration via environment variables only

## 8. Assumptions

- AWS credentials available via environment or IAM role
- PostgreSQL instance accessible from application
- S3 bucket exists and is properly configured
- GitHub repository has package write permissions

## 9. Dependencies

| Dependency | Type | Description |
|------------|------|-------------|
| PostgreSQL 16+ | External | Database service |
| AWS S3 | External | Object storage |
| GitHub Actions | External | CI/CD platform |
| GHCR | External | Container registry |

## 10. Timeline

| Phase | Deliverable |
|-------|-------------|
| Phase 1 | Core application with DB health check |
| Phase 2 | S3 upload functionality |
| Phase 3 | File management CRUD |
| Phase 4 | Docker + CI/CD pipeline |

---

**Approval:** Document approved for implementation.
