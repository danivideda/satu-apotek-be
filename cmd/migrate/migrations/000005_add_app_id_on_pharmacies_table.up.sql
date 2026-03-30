ALTER TABLE IF EXISTS pharmacies
ADD COLUMN IF NOT EXISTS app_id varchar(15) UNIQUE NOT NULL DEFAULT '0';

CREATE INDEX pharmacies_app_id_idx ON pharmacies (app_id);

CREATE TABLE IF NOT EXISTS sqids_config(
  alphabets text PRIMARY KEY
);