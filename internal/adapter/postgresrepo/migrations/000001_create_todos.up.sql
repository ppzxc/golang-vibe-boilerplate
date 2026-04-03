CREATE TABLE IF NOT EXISTS todos (
    id          VARCHAR(36)  PRIMARY KEY,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    completed   BOOLEAN      NOT NULL DEFAULT FALSE,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
