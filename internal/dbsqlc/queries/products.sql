-- name: ProductsGetByID :one
SELECT 
  p.id,
  p.pharmacy_id, 
  p.name, 
  p.sku, 
  p.barcode, 
  p.min_stock, 
  p.base_sell_price, pu.name AS base_unit_name, 
  pu.label AS base_unit_label
FROM products p
LEFT JOIN product_units pu ON pu.id = p.base_unit_id
WHERE p.id = @id;


-- name: ProductsGetByPharmacyID :many
SELECT
  sqlc.embed(p),
  pc.code AS category_code,
  pc.name AS category_name,
  pu.label AS unit_label,
  pu.name AS unit_name
FROM products p
LEFT JOIN product_categories pc ON pc.id = p.category_id
LEFT JOIN product_units pu ON pu.id = p.base_unit_id
WHERE p.pharmacy_id = @pharmacy_id
ORDER BY p.id ASC
LIMIT sqlc.arg('limit') OFFSET sqlc.arg('offset');

-- name: ProductsGetTotalByPharmacyID :one
SELECT COUNT(*) FROM products WHERE pharmacy_id = @pharmacy_id;

-- name: ProductsCreate :one
INSERT INTO products (pharmacy_id, "name", sku, barcode, category_id, base_unit_id, min_stock, base_sell_price, is_active)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING *;