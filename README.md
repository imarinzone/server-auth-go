# Server-to-Server Auth Service

A robust, production-ready authentication service written in Go. It implements the **OAuth 2.0 Client Credentials Flow** to issue JWTs for server-to-server communication.

## Features

- **JWT Authentication**: Secure stateless authentication using JSON Web Tokens (HS256).
- **Pluggable Storage**:
    - **PostgreSQL**: Persistent storage for client credentials.
    - **Redis**: Caching layer for high-performance credential verification.
    - **In-Memory**: Fallback for local development.
- **Dockerized**: Full stack containerization with Docker Compose.
- **Middleware**: Reusable Go middleware for protecting downstream services.

## Prerequisites

- Go 1.23+
- Docker & Docker Compose (optional, for full stack)

## Getting Started

### Option 1: Docker (Recommended)

Run the entire stack (Auth Service, Postgres, Redis):

```bash
docker-compose up --build
```

The service will be available at `http://localhost:8080`.

### Option 2: Local Development

1.  **Install Dependencies**:
    ```bash
    go mod download
    ```

2.  **Run the Server**:
    ```bash
    # Runs with in-memory store by default if DB env vars are missing
    go run main.go
    ```

## API Usage

### 1. Get Access Token

Exchange your Client ID and Secret for a JWT.

**Request:**
```bash
POST /token
Content-Type: application/json

{
  "client_id": "service-a",
  "client_secret": "secret-a"
}
```

**Response:**
```json
{
  "access_token": "eyJhbGciOiJIUzI1Ni...",
  "expires_in": 3600,
  "token_type": "Bearer"
}
```

### 2. Access Protected Resource

Use the token to access protected endpoints.

**Request:**
```bash
GET /protected
Authorization: Bearer <YOUR_ACCESS_TOKEN>
```

## Configuration

The application is configured via environment variables:

| Variable | Description | Default |
| :--- | :--- | :--- |
| `PORT` | Server port | `8080` |
| `SECRET_KEY` | JWT signing key | `super-secret...` |
| `ISSUER` | JWT issuer claim | `auth-server` |
| `POSTGRES_CONN` | Postgres connection string | `host=localhost...` |
| `REDIS_ADDR` | Redis address | `localhost:6379` |

## Project Structure

- `cmd/`: Application entry points.
- `internal/auth/`: Authentication logic and storage adapters (Postgres, Redis, Memory).
- `internal/token/`: JWT generation and validation.
- `pkg/middleware/`: Reusable HTTP middleware for other services.
