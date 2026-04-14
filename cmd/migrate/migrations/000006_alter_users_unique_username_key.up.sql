BEGIN;
  -- Drop the old unique username constraint
  ALTER TABLE users DROP CONSTRAINT IF EXISTS users_username_key;
  -- Add the new composite unique key
  ALTER TABLE users ADD CONSTRAINT users_username_pharmacy_id_unique_key UNIQUE (username, pharmacy_id);
COMMIT;