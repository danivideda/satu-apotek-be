-- name: GetProductsByPharmacyID :many
SELECT * FROM products
WHERE pharmacy_id = @pharmacy_id;

-- name: CreateProduct :one
INSERT INTO products (pharmacy_id, "name", min_stock, base_sell_price)
VALUES ($1, $2, $3, $4) RETURNING *;