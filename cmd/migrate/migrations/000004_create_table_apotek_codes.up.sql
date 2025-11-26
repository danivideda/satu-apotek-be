CREATE TABLE IF NOT EXISTS apotek_codes (
    apotek_id uuid PRIMARY KEY,
    code text UNIQUE NOT NULL,
    expires timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW(),
    FOREIGN KEY (apotek_id) REFERENCES pharmacies (id) 
);