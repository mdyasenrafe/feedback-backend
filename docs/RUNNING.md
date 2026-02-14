# RUNNING.md

---

## Prerequisites

| Requirement    | Version / Notes                           |
| -------------- | ----------------------------------------- |
| **Go**         | ≥ 1.25.0 (see `go.mod`)                   |
| **PostgreSQL** | ≥ 12 (local install, Homebrew, or Docker) |
| **psql**       | Needed to run migration SQL files         |

---

## 1. Clone & Install Dependencies

```bash
git clone <repository-url>
cd feedback-backend
go mod download
```

---

## 2. Create `.env`

Copy the example and fill in real values:

```bash
cp .env.example .env
```

Edit `.env` to look like this (**replace every placeholder**):

```bash
# Server
PORT=8080

# Database (local example)
DATABASE_URL=postgres://postgres:password@localhost:5432/feedbackapp?sslmode=disable

# JWT — use a strong random string
JWT_SECRET=CHANGE_ME_TO_A_RANDOM_SECRET

# Deep link — your backend endpoint (or ngrok tunnel for local mobile testing)
APP_DEEPLINK_URL=http://localhost:8080/auth/deeplink

# Mailgun
MAILGUN_API_KEY=key-XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
MAILGUN_DOMAIN=sandboxXXXXXXXXXXXX.mailgun.org
MAILGUN_BASE_URL=https://api.mailgun.net
EMAIL_FROM=FeedbackApp <postmaster@sandboxXXXXXXXXXXXX.mailgun.org>
```

> **Warning:** `.env` is in `.gitignore` — never commit it.

> **Note:** The checked-in `.env.example` references `RESEND_API_KEY` and `EMAIL_FROM` but the actual code (`internal/config/config.go`) reads Mailgun variables. Use the list above.

---

## 3. Load Environment Variables

The application reads environment variables via `os.Getenv`

You **must** export the variables into your shell before running the server:

```bash
set -a
source .env
set +a
```

> `set -a` causes every variable assignment to be exported automatically; `source .env` loads the file; `set +a` disables auto-export.

---

## 4. Set Up PostgreSQL

### Option A — Docker (quick start)

```bash
docker run -d \
  --name feedback-pg \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=feedbackapp \
  -p 5432:5432 \
  postgres:16
```

### Option B — Local Install (macOS)

```bash
brew install postgresql@16
brew services start postgresql@16
createdb feedbackapp
```

---

## 5. Run Migrations

Migration files live in `internal/db/migrations/`.

```bash
psql "$DATABASE_URL" -f internal/db/migrations/001_auth.sql
psql "$DATABASE_URL" -f internal/db/migrations/002_feedback.sql
```

Verify the tables exist:

```bash
psql "$DATABASE_URL" -c "\dt"
```

Expected output:

```
 Schema |    Name     | Type  |  Owner
--------+-------------+-------+----------
 public | feedback    | table | postgres
 public | login_links | table | postgres
 public | users       | table | postgres
```

---

## 6. Start the Server

```bash
go run cmd/api/main.go
```

Expected log output:

```
Database connection established
Server starting on :8080
```

Verified: entrypoint is `cmd/api/main.go` (see line 19 — `func main()`).

---

## 7. Verify the Server Is Up

The server registers a health endpoint at `/health` (`cmd/api/main.go:40-43`):

```bash
curl -s http://localhost:8080/health
```

Expected response:

```
ok
```

---

## 8. (Optional) Hot Reload with Air

An `.air.toml` is included in the repository. It builds `./cmd/api` and outputs to `./tmp/main`.

### Install Air

```bash
go install github.com/air-verse/air@latest
```

### Run

```bash
# Ensure env vars are loaded first
set -a && source .env && set +a

air
```

Air watches for `.go` file changes under `cmd/` and `internal/`, rebuilds, and restarts automatically.

### `.air.toml` (already in repo)

```toml
root = "."
tmp_dir = "tmp"

[build]
  cmd = "go build -o ./tmp/main ./cmd/api"
  bin = "tmp/main"
  delay = 200
  include_ext = ["go"]
  exclude_dir = ["tmp", "vendor"]

[log]
  time = true
```

---

## Stopping the Server

Press `Ctrl+C`. The server performs a **graceful shutdown** within 10 seconds (see `cmd/api/main.go:75-86`).

---

## Troubleshooting

| Symptom                           | Cause                                        | Fix                                               |
| --------------------------------- | -------------------------------------------- | ------------------------------------------------- |
| `panic: DATABASE_URL is required` | Env vars not loaded                          | Run `set -a; source .env; set +a` before starting |
| `failed to ping database`         | Postgres not running or wrong `DATABASE_URL` | Verify with `psql "$DATABASE_URL" -c "SELECT 1"`  |
| `mailgun send failed: status=401` | Invalid `MAILGUN_API_KEY`                    | Verify in Mailgun dashboard                       |
| Port already in use               | Another process on 8080                      | Change `PORT` in `.env` or kill the other process |
