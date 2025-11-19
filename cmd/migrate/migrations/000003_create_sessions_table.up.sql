CREATE TABLE IF NOT EXISTS revoked_sessions (
    session_id text PRIMARY KEY,
    expires timestamp(0) with time zone NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);