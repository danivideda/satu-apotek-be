-- name: CreatePharmacySession :one
INSERT INTO pharmacy_sessions (pharmacy_id, expires_at) 
VALUES ($1, $2) RETURNING *;

-- name: UpdatePharmacySession :one
UPDATE pharmacy_sessions
SET id = gen_random_uuid(),
  expires_at = @expires_at,
  updated_at = @updated_at
WHERE id = @session_id RETURNING *;

-- name: GetPharmacySession :one
SELECT * FROM pharmacy_sessions
WHERE id = @session_id;

-- name: DeletePharmacySession :one
DELETE FROM pharmacy_sessions
WHERE id = @id RETURNING *;

-- name: DeleteExpiredPharmacySessions :many
DELETE FROM pharmacy_sessions
WHERE expires_at < NOW() RETURNING *;