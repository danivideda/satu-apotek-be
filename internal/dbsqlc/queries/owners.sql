-- name: ListOwners :many
SELECT * FROM owners
ORDER BY id;

-- name: CreateOwner :one
INSERT INTO owners (username, email, password_hash) 
VALUES ($1, $2, $3) RETURNING id, email, username, created_at, updated_at;

-- name: GetOwnerByID :one
SELECT * FROM owners
WHERE id = @id;

-- name: GetOwnerByUsername :one
SELECT id, username, password_hash FROM owners
WHERE username = @username;