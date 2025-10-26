-- +migrate Up
CREATE TABLE IF NOT EXISTS qr_codes (
    id VARCHAR(20) PRIMARY KEY,
    text TEXT NOT NULL,
    format VARCHAR(10) NOT NULL DEFAULT 'png',
    size INT NOT NULL DEFAULT 256,
    image_data BYTEA,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_qr_codes_created_at ON qr_codes(created_at DESC);

-- +migrate Down
DROP TABLE IF EXISTS qr_codes;
