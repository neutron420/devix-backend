# Devix Backend - Engineering Knowledge Sharing Platform

Devix is a production-grade, modular monolith backend architected for a developer-centric knowledge sharing and collaboration platform. The system is built with Go, leveraging GORM for schema management and PostgreSQL for persistent storage, with a focus on modularity, security, and real-time performance.

## Architectural Design

The application implements a Modular Monolith architecture following Clean Architecture principles. Each domain is isolated into independent modules to ensure a strict separation of concerns and to facilitate future migration to microservices if required.

### Core Technologies
- Language: Go (Golang) 1.21+
- Framework: Gin Gonic (HTTP Routing)
- Database: PostgreSQL via NeonDB
- ORM: GORM (Automated Schema Synchronization)
- Real-time Engine: Gorilla WebSocket (Event Broadcasting)
- Object Storage: Cloudflare R2 (S3-Compatible)
- Authentication: JWT with Refresh Token Rotation and Argon2id Hashing

---

## Technical Features

### Authentication and Security
The system utilizes a secure JWT-based authentication mechanism. It supports access and refresh token pairs with a rotation policy to mitigate session theft. Password security is enforced using the Argon2id hashing algorithm, providing high resistance to GPU-based brute-force attacks.

### Database Management
Schema management is automated through GORM's Auto-Migration feature, ensuring that the database state always reflects the Go model definitions. The repository layer utilizes fluent API calls for data persistence while maintaining support for complex PostgreSQL features like full-text search.

### Real-time Communication
A centralized WebSocket Hub manages concurrent client connections and broadcasts events (such as new comments or notifications) across the platform.

### Media Management
A flexible storage abstraction layer supports both local filesystem storage for development and Cloudflare R2 for production environments.

---

## Setup and Installation

### Prerequisites
- Go 1.21 or higher
- Access to a PostgreSQL instance
- Cloudflare R2 credentials (for media features)

### Configuration
1. Initialize the environment configuration:
   ```powershell
   cp .env.example .env
   ```
2. Update the `.env` file with the appropriate database URI, JWT secrets, and storage credentials.

### Execution
Run the application using the following command:
```powershell
go run cmd/server/main.go
```
The database schema will be synchronized automatically upon application startup.

---

## Development Roadmap

The backend is designed for iterative expansion. Future implementations include:

1. Advanced Search: Integration of Elasticsearch for high-performance technical search and indexing.
2. Global Caching: Redis implementation for high-traffic data points and session management.
3. Social Authentication: Support for OAuth2 providers, specifically GitHub and Google.
4. Background Processing: Dedicated worker queues for asynchronous tasks and media optimization.
5. Observability: Integration of structured logging and performance monitoring tools.

---

## Project Structure
- cmd/server: Application entry point and dependency injection.
- internal/config: Environment and application configuration.
- internal/database: Database connection and pool management.
- internal/models: Domain model definitions (Source of truth for schema).
- internal/modules: Independent domain modules (Auth, Post, User, etc.).
- internal/pkg: Shared utilities and core library wrappers.

---
Copyright (c) 2026 Devix Engineering. All rights reserved.