CREATE TABLE IF NOT EXISTS feedback (
  id BIGSERIAL PRIMARY KEY,
  user_id BIGINT NOT NULL,
  message TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_feedback_user_id ON feedback(user_id);
CREATE INDEX IF NOT EXISTS idx_feedback_created_at ON feedback(created_at);

-- If your users table uses BIGINT ids:
ALTER TABLE feedback
  ADD CONSTRAINT fk_feedback_user
  FOREIGN KEY (user_id) REFERENCES users(id)
  ON DELETE CASCADE;
