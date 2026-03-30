ALTER TABLE IF EXISTS pharmacies
DROP COLUMN IF EXISTS app_id;

DROP INDEX IF EXISTS pharmacies_app_id_idx;

DROP TABLE IF EXISTS sqids_config;