# API Documentation

## Base URL

Local development: `http://localhost:8080`

## Endpoints

### Health Check

**GET** `/health`

Returns server health status.

**Response:**

- **200 OK**: `ok` (plain text)

---

### Request Login Link

**POST** `/auth/login-link`

Requests a magic link to be sent to the user's email. Creates a new user if the email doesn't exist.

**Request Body:**

```json
{
  "email": "user@example.com"
}
```

**Success Response:**

- **200 OK**

```json
{
  "ok": true
}
```

**Error Responses:**

- **400 Bad Request** - Invalid JSON

```json
{
  "error": "invalid_json"
}
```

- **500 Internal Server Error** - Email send failed

```json
{
  "error": "email_send_failed"
}
```

- **500 Internal Server Error** - Other server errors

```json
{
  "error": "internal_error"
}
```

**Behavior:**

- Email is normalized (trimmed and lowercased) before processing
- User is created if they don't exist (upsert)
- Token expires in 15 minutes
- Response does NOT leak whether the user already existed
- Raw token is NEVER logged or returned in the API response

**Email Format:**

The email contains a deep link in this format:

```
{APP_DEEPLINK_URL}?token=<RAW_TOKEN>
```

Example:

```
myapp://auth?token=AbCdEfGhIjKlMnOpQrStUvWxYz0123456789_-
```

---

### Verify Login Link

**POST** `/auth/login-link/verify`

Verifies a login token and returns a JWT access token.

**Request Body:**

```json
{
  "token": "AbCdEfGhIjKlMnOpQrStUvWxYz0123456789_-"
}
```

**Success Response:**

- **200 OK**

```json
{
  "accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6InVzZXJAZXhhbXBsZS5jb20iLCJzdWIiOiI1NTBlODQwMC1lMjliLTQxZDQtYTcxNi00NDY2NTU0NDAwMDAiLCJpYXQiOjE3MDk1NjgwMDAsImV4cCI6MTcxMjE2MDAwMH0.signature",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "user@example.com"
  }
}
```

**Error Responses:**

- **400 Bad Request** - Invalid JSON

```json
{
  "error": "invalid_json"
}
```

- **401 Unauthorized** - Invalid, expired, or already-used token

```json
{
  "error": "invalid_or_expired_token"
}
```

- **500 Internal Server Error** - Other server errors

```json
{
  "error": "internal_error"
}
```

**Behavior:**

- Token verification is atomic (race-safe)
- Each token can only be used once
- Expired tokens (>15 minutes old) are rejected
- JWT expires in 30 days

**JWT Claims:**

The returned JWT contains:

- `sub`: User ID (UUID string)
- `email`: User email
- `iat`: Issued at (Unix timestamp)
- `exp`: Expires at (Unix timestamp, 30 days from `iat`)

---

## Deep Link Format

The React Native app should register a deep link handler for the scheme defined in `APP_DEEPLINK_URL`.

**Example Configuration:**

If `APP_DEEPLINK_URL=myapp://auth`, the app should handle URLs like:

```
myapp://auth?token=<TOKEN>
```

**React Native Implementation:**

1. Extract the `token` query parameter from the deep link
2. Call `POST /auth/login-link/verify` with the token
3. Store the returned `accessToken` in secure storage (e.g., `@react-native-async-storage/async-storage` or `react-native-keychain`)
4. Use the token in the `Authorization: Bearer <accessToken>` header for future API requests

---

## Error Codes

| Code                       | HTTP Status | Meaning                                    |
| -------------------------- | ----------- | ------------------------------------------ |
| `invalid_json`             | 400         | Request body is not valid JSON             |
| `invalid_or_expired_token` | 401         | Token is invalid, expired, or already used |
| `email_send_failed`        | 500         | Failed to send email via SMTP              |
| `internal_error`           | 500         | Unexpected server error                    |
| `method_not_allowed`       | 405         | HTTP method not allowed for this endpoint  |

---

## Security Notes

- **Token Storage**: Only SHA256 hashes of tokens are stored in the database
- **Token Expiry**: Login tokens expire in 15 minutes
- **JWT Expiry**: Access tokens expire in 30 days
- **Race Conditions**: Token consumption uses atomic SQL (`UPDATE...RETURNING`) to prevent double-use
- **Email Normalization**: Emails are stored lowercase to prevent duplicate accounts
- **No Token Logging**: Raw tokens are never logged on the server
