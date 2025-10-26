-- +migrate Up
CREATE TABLE IF NOT EXISTS pastes (
    id VARCHAR(20) PRIMARY KEY,
    title VARCHAR(255),
    content TEXT NOT NULL,
    syntax VARCHAR(50),
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    is_compressed BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_pastes_created_at ON pastes(created_at DESC);
CREATE INDEX idx_pastes_expires_at ON pastes(expires_at);
CREATE INDEX idx_pastes_is_public ON pastes(is_public);

-- +migrate Down
DROP TABLE IF EXISTS pastes;
