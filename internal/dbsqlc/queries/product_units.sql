-- name: ProductUnitsGetByPharmacyID :many
SELECT * FROM product_units
WHERE pharmacy_id = @pharmacy_id;


-- name: ProductUnitsCreate :one
INSERT INTO product_units (pharmacy_id, "name", label)
VALUES ($1, $2, $3) RETURNING *;