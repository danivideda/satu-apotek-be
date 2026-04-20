-- name: CreateUserSession :one
INSERT INTO user_sessions (user_id, expires_at) 
VALUES ($1, $2) RETURNING *;

-- name: UpdateUserSession :one
UPDATE user_sessions
SET id = gen_random_uuid(),
  expires_at = @expires_at,
  updated_at = @updated_at
WHERE id = @session_id RETURNING *;

-- name: GetUserSession :one
SELECT * FROM user_sessions
WHERE id = @session_id;

-- name: DeleteUserSession :one
DELETE FROM user_sessions
WHERE id = @id RETURNING *;

-- name: DeleteExpiredUserSessions :many
DELETE FROM user_sessions
WHERE expires_at < NOW() RETURNING *;