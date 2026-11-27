-- name: ProductUnitConversionsCreate :one
INSERT INTO product_unit_conversions (product_id, unit_id, quantity_in_base_unit)
VALUES ($1, $2, $3) RETURNING *;