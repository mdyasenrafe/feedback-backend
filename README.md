# README.md

This is a simple backend API for a feedback application. It provides passwordless authentication via email magic links and allows authenticated users to submit feedback.

## Table of Contents

- [Project Overview](#project-overview)
- [Tech Stack](#tech-stack)
- [Feature-Based Folder Structure](#feature-based-folder-structure)
- [Configuration](#configuration)
  - [How Environment Variables Are Loaded](#how-environment-variables-are-loaded)
  - [Required Environment Variables](#required-environment-variables)
  - [.env.example](#envexample)
- [Postman Collection](#postman-collection)
  - [Collection Variables](#collection-variables)
  - [Auth Workflow](#auth-workflow)
- [Documentation](#documentation)

## Project Overview

**FeedbackApp Backend** is a Go HTTP API server that provides:

- **Passwordless authentication** via email magic links (Mailgun).
- **Feedback collection** from authenticated users, with optional Slack publishing (currently a mock implementation).

The server exposes a small, focused API:

| Method | Path                      | Auth?   | Purpose                                       |
| ------ | ------------------------- | ------- | --------------------------------------------- |
| GET    | `/health`                 | No      | Liveness / readiness probe                    |
| POST   | `/auth/login-link`        | No      | Send a magic-link email to the user           |
| POST   | `/auth/login-link/verify` | No      | Exchange the magic-link token for a JWT       |
| GET    | `/auth/deeplink`          | No      | HTML page that opens the mobile app deep link |
| POST   | `/feedback`               | **Yes** | Submit feedback (requires Bearer JWT)         |

For full endpoint details see [docs/API.md](docs/API.md).

---

## Tech Stack

| Component    | Detail                                                                |
| ------------ | --------------------------------------------------------------------- |
| **Language** | Go 1.25.0                                                             |
| **HTTP**     | Standard library `net/http` with `http.ServeMux`                      |
| **Database** | PostgreSQL via `github.com/jackc/pgx/v5` (connection pool: `pgxpool`) |

---

## Feature-Based Folder Structure

```
feedback-backend/
├── cmd/
│   └── api/
│       └── main.go                    # Entrypoint — server bootstrap, route registration
├── internal/
│   ├── config/
│   │   └── config.go                  # Reads env vars via os.Getenv; panics on missing required vars
│   ├── db/
│   │   ├── db.go                      # pgxpool.Pool initialisation + ping
│   │   └── migrations/
│   │       ├── 001_auth.sql           # DDL: users, login_links tables + indexes
│   │       └── 002_feedback.sql       # DDL: feedback table + indexes
│   ├── middleware/
│   │   └── auth.go                    # JWT Bearer token validation middleware
│   ├── modules/
│   │   ├── auth/                      # Authentication module
│   │   │   ├── auth.handler.go        # HTTP handlers (login-link, verify, deeplink)
│   │   │   ├── auth.service.go        # Business logic (request link, verify link)
│   │   │   ├── auth.repo.go           # Database queries (upsert user, create/consume link)
│   │   │   ├── auth.jwt.go            # JWT creation (HS256, 10-year expiry)
│   │   │   ├── auth.tokens.go         # Secure random token generation + SHA-256 hashing
│   │   │   ├── auth.mail.go           # Mailgun email client
│   │   │   ├── auth.routes.go         # Route registration on ServeMux
│   │   │   └── auth.types.go          # Request/Response/Domain structs
│   │   └── feedback/                  # Feedback module
│   │       ├── feedback.handler.go    # HTTP handler (create feedback)
│   │       ├── feedback.service.go    # Business logic (validate, persist, publish)
│   │       ├── feedback.repo.go       # Database queries (INSERT … RETURNING)
│   │       ├── feedback.routes.go     # Route registration with auth middleware
│   │       ├── feedback.types.go      # Request/Response/Domain structs
│   │       ├── slack.go               # SlackClient interface
│   │       └── slack.mock.go          # Mock implementation (logs to stdout)
│   └── shared/
│       └── httpx/
│           └── json.go                # WriteJSON / WriteError helpers
├── .air.toml                          # Air hot-reload config
├── .env.example                       # Template for environment variables
├── .gitignore                         # Ignores .env
├── go.mod
└── go.sum
```

**Design:** Each module (`auth`, `feedback`) is self-contained with its own handler → service → repository layers. Modules only depend on `shared/httpx` and `middleware`, never on each other.

---

## Configuration

### How Environment Variables Are Loaded

The application uses `os.Getenv` exclusively (see `internal/config/config.go`). You must export the variables into your shell **before** starting the server:

```bash
set -a          # auto-export every variable sourced
source .env     # load your .env file
set +a          # stop auto-exporting
```

### Required Environment Variables

Derived from `internal/config/config.go`:

| Variable           | Required? | Default                   | Purpose                                                                |
| ------------------ | --------- | ------------------------- | ---------------------------------------------------------------------- |
| `PORT`             | No        | `8080`                    | HTTP listen port                                                       |
| `DATABASE_URL`     | **Yes**   | —                         | PostgreSQL connection string                                           |
| `JWT_SECRET`       | **Yes**   | —                         | HMAC-SHA256 key for signing JWTs                                       |
| `APP_DEEPLINK_URL` | **Yes**   | —                         | Base URL of the `/auth/deeplink` endpoint (backend appends `?token=…`) |
| `MAILGUN_API_KEY`  | **Yes**   | —                         | Mailgun API key                                                        |
| `MAILGUN_DOMAIN`   | **Yes**   | —                         | Mailgun sending domain (e.g. `sandbox…mailgun.org`)                    |
| `MAILGUN_BASE_URL` | No        | `https://api.mailgun.net` | Mailgun API base (use `https://api.eu.mailgun.net` for EU)             |
| `EMAIL_FROM`       | **Yes**   | —                         | Sender address (e.g. `FeedbackApp <postmaster@sandbox…mailgun.org>`)   |

### `.env.example`

```bash
# ── Server ────────────────────────────────────────
PORT=8080

# ── Database ──────────────────────────────────────
DATABASE_URL=postgres://postgres:password@localhost:5432/feedbackapp?sslmode=disable

# ── JWT ───────────────────────────────────────────
JWT_SECRET=CHANGE_ME_TO_A_RANDOM_SECRET

# ── Deep link ─────────────────────────────────────
# Points to the backend /auth/deeplink endpoint (or tunnel URL during local dev)
APP_DEEPLINK_URL=http://localhost:8080/auth/deeplink

# ── Mailgun ───────────────────────────────────────
MAILGUN_API_KEY=key-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
MAILGUN_DOMAIN=sandboxXXXXXXXXXXXX.mailgun.org
MAILGUN_BASE_URL=https://api.mailgun.net
EMAIL_FROM=FeedbackApp <postmaster@sandboxXXXXXXXXXXXX.mailgun.org>
```

---

## Postman Collection

Use this collection to run the backend API requests from Postman (local dev and deployed). It’s a quick way to validate the auth + feedback flow without writing curl commands.

[<img src="https://run.pstmn.io/button.svg" alt="Run In Postman" style="width: 128px; height: 32px;">](https://god.gw.postman.com/run-collection/36553719-b4a49e44-4636-494c-bb99-323664e5c49d?action=collection%2Ffork&source=rip_markdown&collection-url=entityId%3D36553719-b4a49e44-4636-494c-bb99-323664e5c49d%26entityType%3Dcollection%26workspaceId%3Dc9b500b3-e7af-43ae-918a-f7ff71340372)

**How to use**

1. Click **Run in Postman** to fork the collection into your workspace.
2. Set the base URL to match where the server is running:
   - Local: `http://localhost:8080`
   - Deployed: your hosted API URL (e.g., Render)
3. Run the requests in order (start with the health check, then auth, then feedback).

**Notes**

- If the collection includes authenticated requests, make sure you paste the returned access token into the token variable (or Authorization tab) before calling protected endpoints.
- The collection is intended to match the repo docs (`API.md`) and the backend route definitions—if an endpoint changes, update the collection accordingly.

### Collection Variables

| Variable      | Default Value                | Purpose                                         |
| ------------- | ---------------------------- | ----------------------------------------------- |
| `localUrl`    | `http://localhost:8080`      | Base URL for local development                  |
| `baseUrl`     | `https://YOUR-PROD-BASE-URL` | Base URL for production — replace with your URL |
| `bearerToken` | _(empty)_                    | JWT access token — set after login              |

All requests default to `{{localUrl}}`. To hit production, change the URL prefix to `{{baseUrl}}`.

### Auth Workflow

1. Send **Request Login Link** (`POST /auth/login-link`) to trigger a magic-link email.
2. Obtain the raw token from the email link.
3. Send **Verify Login Link** (`POST /auth/login-link/verify`) with the token.
4. Copy the `accessToken` value from the response.
5. Set the `bearerToken` collection variable to that value.
6. Authenticated requests (e.g. `POST /feedback`) will use `Authorization: Bearer {{bearerToken}}` automatically.

Cross-reference request schemas with [API.md](API.md) for the full specification.

---

## Documentation

- [Running locally](docs/RUNNING.md)
- [API reference](docs/API.md)
- [Database schema](docs/DATABASE_SCHEMA.md)
