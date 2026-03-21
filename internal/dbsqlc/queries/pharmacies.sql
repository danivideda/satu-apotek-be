-- name: CreatePharmacy :one
INSERT INTO pharmacies (owner_id, name) 
VALUES ($1, $2) RETURNING *;

-- name: GetPharmaciesByOwner :many
SELECT * FROM pharmacies
WHERE owner_id = @owner_id
ORDER BY created_at DESC;