# Running the FeedbackApp Auth Backend

This guide covers local development setup and testing. **No Docker required.**

## Prerequisites

- Go 1.22 or later
- PostgreSQL 14+ running locally

## Step 1: Set Up PostgreSQL

### Create Database

```bash
createdb feedbackapp
```

### Apply Migration

```bash
psql feedbackapp < internal/db/migrations/001_auth.sql
```

### Verify Tables

```bash
psql feedbackapp -c "\dt"
```

You should see `users` and `login_links` tables.

## Step 2: Configure Environment

Copy the example environment file:

```bash
cp .env.example .env
```

Edit `.env` with your actual values:

```bash
# Database (update if needed)
DATABASE_URL=postgres://localhost:5432/feedbackapp?sslmode=disable

# JWT Secret (generate a random string)
JWT_SECRET=your-random-secret-here

# Deep Link (use your app's scheme)
APP_DEEPLINK_URL=myapp://auth

# SMTP (example with Gmail)
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=your-email@gmail.com
SMTP_PASS=your-app-password
SMTP_FROM=your-email@gmail.com
```

**Note**: For Gmail, you need to use an [App Password](https://support.google.com/accounts/answer/185833), not your regular password.

## Step 3: Install Dependencies

```bash
go mod download
```

## Step 4: Run the Server

```bash
source .env  # Load environment variables
go run cmd/api/main.go
```

You should see:

```
Database connection established
Server starting on :8080
```

## Step 5: Test the Endpoints

### Health Check

```bash
curl http://localhost:8080/health
```

Expected: `ok`

### Request Login Link

```bash
curl -X POST http://localhost:8080/auth/login-link \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'
```

Expected response:

```json
{ "ok": true }
```

Check your email inbox for the login link. It will look like:

```
myapp://auth?token=<LONG_TOKEN_STRING>
```

### Verify Login Link

Copy the token from the email and verify it:

```bash
curl -X POST http://localhost:8080/auth/login-link/verify \
  -H "Content-Type: application/json" \
  -d '{"token":"<PASTE_TOKEN_HERE>"}'
```

Expected response:

```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "test@example.com"
  }
}
```

### Test Error Cases

**Invalid token:**

```bash
curl -X POST http://localhost:8080/auth/login-link/verify \
  -H "Content-Type: application/json" \
  -d '{"token":"invalid"}'
```

Expected: `{"error":"invalid_or_expired_token"}` (401)

**Double-use (use same token twice):**

Run the verify request again with the same token. Expected: `{"error":"invalid_or_expired_token"}` (401)

**Expired token:**

Wait 15+ minutes after requesting a login link, then try to verify. Expected: `{"error":"invalid_or_expired_token"}` (401)

## Development Tips

### View Database Records

```bash
# View users
psql feedbackapp -c "SELECT * FROM users;"

# View login links
psql feedbackapp -c "SELECT id, user_id, expires_at, used_at, created_at FROM login_links ORDER BY created_at DESC LIMIT 10;"
```

### Reset Database

```bash
psql feedbackapp -c "DROP TABLE IF EXISTS login_links CASCADE; DROP TABLE IF EXISTS users CASCADE;"
psql feedbackapp < internal/db/migrations/001_auth.sql
```

### Build for Production

```bash
go build -o bin/api cmd/api/main.go
./bin/api
```

## Troubleshooting

### "Database connection error"

- Verify PostgreSQL is running: `psql -l`
- Check `DATABASE_URL` in `.env`

### "email_send_failed"

- Verify SMTP credentials
- For Gmail, ensure you're using an App Password
- Check firewall/network settings for port 587

### "Configuration error: X is required"

- Ensure all variables in `.env.example` are set in your `.env`
- Run `source .env` before starting the server
