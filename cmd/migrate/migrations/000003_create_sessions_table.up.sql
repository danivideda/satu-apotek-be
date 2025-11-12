CREATE TABLE IF NOT EXISTS revoked_tokens (
    id serial PRIMARY KEY,
    refresh_token text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);