CREATE TABLE IF NOT EXISTS apotek_codes (
    apotek_id bigint PRIMARY KEY REFERENCES pharmacies,
    code text UNIQUE NOT NULL,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX apotek_codes_expires_at_idx ON apotek_codes (expires_at);