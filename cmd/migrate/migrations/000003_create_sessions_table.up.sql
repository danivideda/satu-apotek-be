CREATE TABLE IF NOT EXISTS owner_sessions (
    id serial PRIMARY KEY,
    refresh_token text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_sessions (
    id serial PRIMARY KEY,
    owner_id int REFERENCES owners(id),
    refresh_token text NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);