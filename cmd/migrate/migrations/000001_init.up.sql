CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS owners(
    id bigserial PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    username text UNIQUE NOT NULL,
    password_hash text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pharmacies(
    id bigserial PRIMARY KEY,
    owner_id bigint NOT NULL REFERENCES owners,
    name text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS users(
    id bigserial PRIMARY KEY,
    username text UNIQUE NOT NULL,
    password_hash text NOT NULL,
    pharmacy_id bigint NOT NULL REFERENCES pharmacies,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    updated_at timestamptz NOT NULL DEFAULT NOW()
);