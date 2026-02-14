# API.md

---

## Base URL

| Environment | URL                                                                                      |
| ----------- | ---------------------------------------------------------------------------------------- |
| Local       | `http://localhost:8080` (default `PORT` in `internal/config/config.go`)                  |
| Production  | [https://feedback-backend-shh8.onrender.com](https://feedback-backend-shh8.onrender.com) |

---

## Error Response Format

All errors use a consistent JSON envelope (see `internal/shared/httpx/json.go:16-18`):

```json
{
  "error": "<error_code>"
}
```

Content-Type: `application/json`

---

## Endpoints

---

### 1 · `GET /health`

**Auth:** None

#### Request

```bash
curl http://localhost:8080/health
```

#### Response

| Status   | Body              |
| -------- | ----------------- |
| `200 OK` | `ok` (plain text) |

---

### 2 · `POST /auth/login-link`

Send a magic-link email to the given address.

**Auth:** None

#### Request

```bash
curl -X POST http://localhost:8080/auth/login-link \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'
```

**Body schema:**

```json
{
  "email": "string (required)"
}
```

#### Success Response — `200 OK`

```json
{
  "ok": true
}
```

#### Error Responses

| Status | Error Code           | Condition                      |
| ------ | -------------------- | ------------------------------ |
| `400`  | `invalid_json`       | Request body is not valid JSON |
| `405`  | `method_not_allowed` | Method is not POST             |
| `500`  | `email_send_failed`  | Mailgun API returned an error  |
| `500`  | `internal_error`     | Other server-side error        |

---

### 3 · `POST /auth/login-link/verify`

Exchange the raw magic-link token for a JWT.

**Auth:** None

#### Request

```bash
curl -X POST http://localhost:8080/auth/login-link/verify \
  -H "Content-Type: application/json" \
  -d '{"token":"<raw_token_from_email>"}'
```

**Body schema:**

```json
{
  "token": "string (required — the raw token from the email link)"
}
```

#### Success Response — `200 OK`

```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIs…",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com"
  }
}
```

#### Error Responses

| Status | Error Code                 | Condition                                 |
| ------ | -------------------------- | ----------------------------------------- |
| `400`  | `invalid_json`             | Request body is not valid JSON            |
| `401`  | `invalid_or_expired_token` | Token not found, expired, or already used |
| `405`  | `method_not_allowed`       | Method is not POST                        |
| `500`  | `internal_error`           | Other server-side error                   |

---

### 4 · `GET /auth/deeplink`

Serves an HTML page that attempts to open the native app via deep link (`feedbackapp://auth?token=…`).

**Auth:** None

#### Request

```bash
curl "http://localhost:8080/auth/deeplink?token=SOME_RAW_TOKEN"
```

**Query parameters:**

| Param   | Required | Description          |
| ------- | -------- | -------------------- |
| `token` | Yes      | Raw magic-link token |

#### Success Response — `200 OK`

Content-Type: `text/html; charset=utf-8`

The HTML page contains:

- A JavaScript redirect to `feedbackapp://auth?token=<url-encoded-token>`.
- A fallback button for in-app browsers that block automatic redirects.
- A note that the link expires in 15 minutes.

#### Error Response

| Status | Body                             | Condition                        |
| ------ | -------------------------------- | -------------------------------- |
| `400`  | `missing token` (plain text)     | `?token=` query parameter absent |
| `405`  | `{"error":"method_not_allowed"}` | Method is not GET                |

> **Note:** This endpoint is typically opened by the user clicking the email link in a mobile browser. It is not called directly by the mobile app. The app intercepts `feedbackapp://auth?token=…` and then calls `POST /auth/login-link/verify`.

---

### 5 · `POST /feedback`

Submit a feedback message. **Requires authentication.**

**Auth:** JWT Bearer token required

#### Request

```bash
curl -X POST http://localhost:8080/feedback \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <accessToken>" \
  -d '{"message":"The app is great!"}'
```

**Body schema:**

```json
{
  "message": "string (required — must not be empty after trimming)"
}
```

#### Success Response — `201 Created`

```json
{
  "id": "660e8400-e29b-41d4-a716-446655440000",
  "user_id": "550e8400-e29b-41d4-a716-446655440000",
  "message": "The app is great!",
  "created_at": "2026-02-14T10:30:00Z"
}
```

#### Error Responses

| Status | Error Code           | Condition                                                                      |
| ------ | -------------------- | ------------------------------------------------------------------------------ |
| `400`  | `invalid_json`       | Request body is not valid JSON                                                 |
| `400`  | `message_required`   | Message is empty or whitespace-only (validated in `feedback.service.go:27-28`) |
| `401`  | _(see Auth section)_ | Missing, malformed, or expired JWT                                             |
| `401`  | `unauthorized`       | Context has no user (should not happen if middleware runs)                     |
| `401`  | `invalid_user_id`    | User ID from JWT is not a valid UUID                                           |
| `405`  | `method_not_allowed` | Method is not POST                                                             |
| `500`  | `internal_error`     | Database or other server error                                                 |

---

## Summary Table

| Method | Path                      | Auth   | Success Status | description           |
| ------ | ------------------------- | ------ | -------------- | --------------------- |
| GET    | `/health`                 | None   | `200`          | Health check          |
| POST   | `/auth/login-link`        | None   | `200`          | Generate a login link |
| POST   | `/auth/login-link/verify` | None   | `200`          | Verify a login link   |
| GET    | `/auth/deeplink`          | None   | `200`          | Deep link to the app  |
| POST   | `/feedback`               | Bearer | `201`          | Submit feedback       |
