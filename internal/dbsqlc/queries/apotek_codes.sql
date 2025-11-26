-- name: GetApotekCode :one
SELECT * FROM apotek_codes
WHERE apotek_id = @apotek_id;

-- name: CreateApotekCode :one
INSERT INTO apotek_codes (apotek_id, code, expires) 
VALUES ($1, $2, $3) RETURNING *;

-- name: GetApotekCodeByCode :one
SELECT * FROM apotek_codes
WHERE code = @code;