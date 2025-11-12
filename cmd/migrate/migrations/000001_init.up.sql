CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS owners(
    id serial PRIMARY KEY,
    email citext UNIQUE NOT NULL,
    username varchar(255) UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pharmacies(
    id serial PRIMARY KEY,
    owner_id INT NOT NULL,
    name VARCHAR(255),
    FOREIGN KEY (owner_id) REFERENCES owners (id)
);

CREATE TABLE IF NOT EXISTS users(
    id serial PRIMARY KEY,
    username VARCHAR(255) UNIQUE NOT NULL,
    password_hash bytea NOT NULL,
    pharmacy_id INT NOT NULL,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    FOREIGN KEY (pharmacy_id) REFERENCES pharmacies(id)
);