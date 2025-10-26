-- +migrate Up
CREATE TABLE IF NOT EXISTS short_urls (
    id BIGSERIAL PRIMARY KEY,
    code VARCHAR(20) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    alias VARCHAR(50),
    clicks BIGINT NOT NULL DEFAULT 0,
    is_public BOOLEAN NOT NULL DEFAULT TRUE,
    expires_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS url_clicks (
    id BIGSERIAL PRIMARY KEY,
    short_url_id BIGINT NOT NULL REFERENCES short_urls(id) ON DELETE CASCADE,
    referrer TEXT,
    user_agent TEXT,
    ip_address VARCHAR(45),
    clicked_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_short_urls_code ON short_urls(code);
CREATE INDEX idx_short_urls_expires_at ON short_urls(expires_at);
CREATE INDEX idx_url_clicks_short_url_id ON url_clicks(short_url_id);

-- +migrate Down
DROP TABLE IF EXISTS url_clicks;
DROP TABLE IF EXISTS short_urls;
