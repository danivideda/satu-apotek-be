CREATE TABLE IF NOT EXISTS owner_sessions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    owner_id bigint NOT NULL REFERENCES owners,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_sessions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id bigint NOT NULL REFERENCES users,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pharmacy_sessions (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    pharmacy_id bigint NOT NULL REFERENCES pharmacies,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);