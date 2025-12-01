-- name: CreateOwnerSession :one
INSERT INTO owner_sessions (id, owner_id, expires_at) 
VALUES ($1, $2, $3) RETURNING *;

-- name: UpdateOwnerSession :one
UPDATE owner_sessions
SET id = gen_random_uuid(),
  expires_at = @expires_at,
  updated_at = @updated_at
WHERE id = @session_id RETURNING *;

-- name: GetOwnerSession :one
SELECT * FROM owner_sessions
WHERE id = @session_id;

-- name: DeleteOwnerSession :one
DELETE FROM owner_sessions
WHERE id = @id RETURNING *;