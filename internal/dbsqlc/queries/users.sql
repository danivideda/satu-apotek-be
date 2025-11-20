-- name: CreateUser :one
INSERT INTO users (username, password_hash, pharmacy_id) 
VALUES ($1, $2, $3) RETURNING *;