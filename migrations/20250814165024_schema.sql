-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS url (
  id SERIAL PRIMARY KEY,
  correlation_id TEXT NOT NULL,
  original TEXT UNIQUE NOT NULL,
  short TEXT UNIQUE NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  is_deleted BOOLEAN DEFAULT FALSE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_short_url ON url (short);
CREATE UNIQUE INDEX IF NOT EXISTS idx_original_url ON url (original);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_original_url;
DROP INDEX IF EXISTS idx_short_url;

DROP TABLE IF EXISTS url;
-- +goose StatementEnd
