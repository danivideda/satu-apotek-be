CREATE TABLE IF NOT EXISTS revoked_refresh_tokens (
    id serial PRIMARY KEY,
    token text UNIQUE NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);