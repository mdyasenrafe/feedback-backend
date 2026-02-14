# DATABASE_SCHEMA.md

> **Source of Truth:** This documentation is derived directly from the migration SQL files in `internal/db/migrations/`.

---

## Overview

| Table         | Migration File     | Purpose                              |
| ------------- | ------------------ | ------------------------------------ |
| `users`       | `001_auth.sql`     | User accounts (email-based identity) |
| `login_links` | `001_auth.sql`     | One-time magic-link tokens           |
| `feedback`    | `002_feedback.sql` | User-submitted feedback messages     |

All primary keys are `UUID` (auto-generated via `gen_random_uuid()`). All timestamps are `TIMESTAMPTZ` (UTC-aware).

---

## Tables

### `users`

**Source:** `internal/db/migrations/001_auth.sql:2-6`

```sql
CREATE TABLE users (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    email      TEXT        UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
```

| Column       | Type          | Constraints              | Notes                                                              |
| ------------ | ------------- | ------------------------ | ------------------------------------------------------------------ |
| `id`         | `UUID`        | PK, auto-generated       | —                                                                  |
| `email`      | `TEXT`        | `UNIQUE NOT NULL`        | Normalised to lowercase by application code (`auth.service.go:30`) |
| `created_at` | `TIMESTAMPTZ` | `NOT NULL DEFAULT now()` | —                                                                  |

**Implicit indexes:** PK index on `id`, unique index on `email`.

**Relationships:**

- `users.id` ← `login_links.user_id` (one-to-many, `ON DELETE CASCADE`)
- `users.id` ← `feedback.user_id` (one-to-many, `ON DELETE CASCADE`)

**Application behaviour:** Users are upserted on each login-link request (`INSERT … ON CONFLICT (email) DO UPDATE` — `auth.repo.go:25-30`). There is no password column — authentication is entirely magic-link-based.

---

### `login_links`

**Source:** `internal/db/migrations/001_auth.sql:9-20`

```sql
CREATE TABLE login_links (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT        UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at    TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_login_links_expires_at ON login_links(expires_at);
CREATE INDEX idx_login_links_token_hash ON login_links(token_hash);
```

| Column       | Type          | Constraints                                       | Notes                                                 |
| ------------ | ------------- | ------------------------------------------------- | ----------------------------------------------------- |
| `id`         | `UUID`        | PK, auto-generated                                | —                                                     |
| `user_id`    | `UUID`        | FK → `users(id)`, `ON DELETE CASCADE`, `NOT NULL` | —                                                     |
| `token_hash` | `TEXT`        | `UNIQUE NOT NULL`                                 | SHA-256 hex of the raw token (`auth.tokens.go:21-24`) |
| `expires_at` | `TIMESTAMPTZ` | `NOT NULL`                                        | 15 minutes from creation (`auth.service.go:51`)       |
| `used_at`    | `TIMESTAMPTZ` | Nullable                                          | `NULL` = unused; set to `now()` on consumption        |
| `created_at` | `TIMESTAMPTZ` | `NOT NULL DEFAULT now()`                          | —                                                     |

**Explicit indexes:**

- `idx_login_links_token_hash` — speeds up token look-up during verification.
- `idx_login_links_expires_at` — supports expiry-based queries / cleanup.

**Security:** Raw tokens are **never stored**; only the SHA-256 hash is persisted. Tokens are **one-time use** — the `ConsumeLoginLink` query atomically checks `used_at IS NULL AND expires_at > now()` and sets `used_at = now()` (`auth.repo.go:55-61`).

---

### `feedback`

**Source:** `internal/db/migrations/002_feedback.sql:2-11`

```sql
CREATE TABLE feedback (
    id         UUID        PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id    UUID        NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    message    TEXT        NOT NULL CHECK (length(trim(message)) > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_feedback_user_id    ON feedback(user_id);
CREATE INDEX idx_feedback_created_at ON feedback(created_at);
```

| Column       | Type          | Constraints                                       | Notes                          |
| ------------ | ------------- | ------------------------------------------------- | ------------------------------ |
| `id`         | `UUID`        | PK, auto-generated                                | —                              |
| `user_id`    | `UUID`        | FK → `users(id)`, `ON DELETE CASCADE`, `NOT NULL` | —                              |
| `message`    | `TEXT`        | `NOT NULL`, `CHECK (length(trim(message)) > 0)`   | No max-length constraint in DB |
| `created_at` | `TIMESTAMPTZ` | `NOT NULL DEFAULT now()`                          | —                              |

**Explicit indexes:**

- `idx_feedback_user_id` — supports per-user lookups.
- `idx_feedback_created_at` — supports chronological sorting/filtering.

---

## Entity-Relationship Diagram

```
┌──────────────────┐
│      users       │
│──────────────────│
│ id          (PK) │◄──────────────┐
│ email       (UQ) │               │
│ created_at       │               │
└──────────────────┘               │
         │                         │
         │ 1:N                     │ 1:N
         ▼                         │
┌──────────────────┐    ┌──────────┴───────┐
│   login_links    │    │    feedback       │
│──────────────────│    │──────────────────│
│ id          (PK) │    │ id          (PK) │
│ user_id     (FK) │    │ user_id     (FK) │
│ token_hash  (UQ) │    │ message          │
│ expires_at       │    │ created_at       │
│ used_at          │    └──────────────────┘
│ created_at       │
└──────────────────┘
```

---

## Running Migrations

There is **no migration runner** built into the application (verified — `internal/db/db.go` only initialises a pool). Apply migrations manually:

```bash
# Ensure DATABASE_URL is set
set -a && source .env && set +a

psql "$DATABASE_URL" -f internal/db/migrations/001_auth.sql
psql "$DATABASE_URL" -f internal/db/migrations/002_feedback.sql
```

Verify:

```bash
psql "$DATABASE_URL" -c "\dt"
```

---

## Full Migration SQL (Appendix)

### `internal/db/migrations/001_auth.sql`

```sql
-- Create users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Create login_links table
CREATE TABLE login_links (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash TEXT UNIQUE NOT NULL,
    expires_at TIMESTAMPTZ NOT NULL,
    used_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Create indexes
CREATE INDEX idx_login_links_expires_at ON login_links(expires_at);
CREATE INDEX idx_login_links_token_hash ON login_links(token_hash);
```

### `internal/db/migrations/002_feedback.sql`

```sql
-- Create feedback table
CREATE TABLE feedback (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  message TEXT NOT NULL CHECK (length(trim(message)) > 0),
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Create indexes
CREATE INDEX idx_feedback_user_id ON feedback(user_id);
CREATE INDEX idx_feedback_created_at ON feedback(created_at);
```
