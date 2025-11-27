-- name: GetApotekCode :one
SELECT * FROM apotek_codes
WHERE apotek_id = @apotek_id;

-- name: CreateApotekCode :one
INSERT INTO apotek_codes (apotek_id, code, expires_at) 
VALUES ($1, $2, $3) RETURNING *;

-- name: GetApotekCodeByCode :one
SELECT * FROM apotek_codes
WHERE code = @code;

-- name: UpsertApotekCode :one
INSERT INTO apotek_codes (apotek_id, code, expires_at) 
VALUES ($1, $2, $3) 
ON CONFLICT (apotek_id) DO UPDATE SET code = excluded.code, expires_at = excluded.expires_at
RETURNING *;

-- name: DeleteExpiredApotekCode :many
DELETE FROM apotek_codes
WHERE expires_at < NOW() RETURNING *;