BEGIN;
  ALTER TABLE users ADD CONSTRAINT users_username_key UNIQUE (username);
  ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_pharmacy_id_unique_key;
COMMIT;