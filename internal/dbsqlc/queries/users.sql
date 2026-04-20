-- name: CreateUser :one
INSERT INTO users (username, password_hash, pharmacy_id) 
VALUES ($1, $2, $3) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users
WHERE id = @id;

-- name: GetUserByUsername :one
SELECT id, username, password_hash FROM users
WHERE username = @username;

-- name: GetUserByPharmacyID :many
SELECT id, username FROM users
WHERE pharmacy_id = @pharmacy_id;