-- name: ProductsGetByPharmacyID :many
SELECT * FROM products
WHERE pharmacy_id = @pharmacy_id;

-- name: ProductsCreate :one
INSERT INTO products (pharmacy_id, "name", sku, barcode, category_id, base_unit_id, min_stock, base_sell_price, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;