CREATE TABLE IF NOT EXISTS revoked_sessions (
    id serial PRIMARY KEY,
    session_id text UNIQUE NOT NULL,
    expires timestamp(0) with time zone NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);