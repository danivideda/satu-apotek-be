CREATE TABLE IF NOT EXISTS pharmacy_codes (
    apotek_id bigint PRIMARY KEY REFERENCES pharmacies,
    code text UNIQUE NOT NULL,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE INDEX pharmacy_codes_expires_idx ON pharmacy_codes (expires_at);