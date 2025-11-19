-- name: CreatePharmacy :one
INSERT INTO pharmacies (owner_id, name) 
VALUES ($1, $2) RETURNING *;