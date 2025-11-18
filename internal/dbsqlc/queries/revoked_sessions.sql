-- name: CreateRevokedSession :one
INSERT INTO revoked_sessions (session_id, expires) 
VALUES ($1, $2) RETURNING *;

-- name: GetRevokedSessionBySessionID :one
SELECT * FROM revoked_sessions
WHERE session_id = @session_id;
