# FeedbackApp Auth Backend

Minimal authentication backend for the FeedbackApp React Native application. This backend handles **authentication only** using email-based magic links and JWT sessions.

## What This Does

This is a production-ready Go backend that implements passwordless authentication:

1. **Request Login Link**: User enters their email in the app
2. **Email Deep Link**: Server generates a one-time token and emails a deep link
3. **Verify Token**: App opens the deep link, extracts the token, and verifies it
4. **JWT Session**: Server returns a JWT access token for authenticated API requests

## Tech Stack

- **Go** with `net/http` (no frameworks)
- **PostgreSQL** for user and token storage
- **JWT (HS256)** for session tokens
- **SMTP** for email delivery via `gomail`

## Security Features

- Tokens are cryptographically random (32 bytes)
- Only token hashes are stored in the database (SHA256)
- Tokens expire in 15 minutes
- Atomic token consumption prevents race conditions
- JWT tokens expire in 30 days
- Email addresses are normalized and stored lowercase

## API Endpoints

- `GET /health` - Health check
- `POST /auth/login-link` - Request a login link
- `POST /auth/login-link/verify` - Verify token and get JWT

See [API.md](API.md) for detailed endpoint specifications.

## Getting Started

See [RUNNING.md](RUNNING.md) for local setup and testing instructions.

## Project Structure

```
feedback-backend/
├── cmd/api/main.go                    # Application entry point
├── internal/
│   ├── config/config.go               # Environment configuration
│   ├── db/
│   │   ├── db.go                      # Database connection
│   │   └── migrations/001_auth.sql    # Auth schema migration
│   ├── shared/httpx/json.go           # HTTP utilities
│   └── modules/auth/                  # Auth module
│       ├── auth.types.go              # Request/response types
│       ├── auth.repo.go               # Database operations
│       ├── auth.tokens.go             # Token generation/hashing
│       ├── auth.jwt.go                # JWT creation
│       ├── auth.mail.go               # Email sending
│       ├── auth.service.go            # Business logic
│       ├── auth.handler.go            # HTTP handlers
│       └── auth.routes.go             # Route registration
└── .env.example                       # Environment variables template
```
