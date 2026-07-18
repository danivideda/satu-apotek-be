-- name: GetPharmacyByIDForOwner :one
SELECT * FROM pharmacies
WHERE id = @id AND owner_id = @owner_id;

-- name: GetPharmacyByAppIDForOwner :one
SELECT * FROM pharmacies
WHERE app_id = @app_id AND owner_id = @owner_id;

-- name: CreatePharmacy :one
INSERT INTO pharmacies (owner_id, name, app_id) 
VALUES ($1, $2, $3) RETURNING *;

-- name: InsertAppID :one
UPDATE pharmacies
SET app_id = @app_id
WHERE id = @id RETURNING *;

-- name: GetPharmaciesByOwner :many
SELECT * FROM pharmacies
WHERE owner_id = @owner_id
ORDER BY created_at DESC;