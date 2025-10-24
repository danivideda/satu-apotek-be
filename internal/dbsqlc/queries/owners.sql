-- name: ListOwners :many
SELECT * FROM owners
ORDER BY id;

-- name: CreateOwner :one
INSERT INTO owners (username, email, password) 
VALUES ($1, $2, $3) RETURNING id, email, username, created_at, updated_at;

-- name: GetOwnerByID :one
SELECT * FROM owners
WHERE id = @id
ORDER BY id;
