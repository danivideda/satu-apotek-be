-- name: CreateRevokedSession :one
INSERT INTO revoked_sessions (session_id, expires) 
VALUES ($1, $2) RETURNING *;