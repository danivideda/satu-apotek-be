-- Drop indexes that reference pharmacy_id
DROP INDEX IF EXISTS product_unit_conversions_pharmacy_idx;
DROP INDEX IF EXISTS product_unit_conversions_pharmacy_and_product_idx;

-- Drop the unique constraint that includes pharmacy_id
ALTER TABLE product_unit_conversions
  DROP CONSTRAINT IF EXISTS product_unit_conversions_pharmacy_id_product_id_unit_id_key;

-- Drop the redundant column (FK is dropped with it)
ALTER TABLE product_unit_conversions
  DROP COLUMN pharmacy_id;

-- Recreate unique constraint scoped only to product + unit
-- (product_id and unit_id are already NOT NULL from 000010)
ALTER TABLE product_unit_conversions
  ADD CONSTRAINT product_unit_conversions_product_id_unit_id_key
  UNIQUE (product_id, unit_id);

-- Useful index for lookups by product
CREATE INDEX product_unit_conversions_product_idx
  ON product_unit_conversions (product_id);