-- Reverse the new unique constraint and index
ALTER TABLE product_unit_conversions
  DROP CONSTRAINT IF EXISTS product_unit_conversions_product_id_unit_id_key;

DROP INDEX IF EXISTS product_unit_conversions_product_idx;

-- Re-add the column (nullable so we can backfill)
ALTER TABLE product_unit_conversions
  ADD COLUMN pharmacy_id bigint;

-- Backfill from products (product_id is already NOT NULL)
UPDATE product_unit_conversions puc
SET pharmacy_id = p.pharmacy_id
FROM products p
WHERE puc.product_id = p.id;

-- Now safe to restore NOT NULL
ALTER TABLE product_unit_conversions
  ALTER COLUMN pharmacy_id SET NOT NULL;

-- Restore the original unique constraint and indexes
ALTER TABLE product_unit_conversions
  ADD CONSTRAINT product_unit_conversions_pharmacy_id_product_id_unit_id_key
  UNIQUE (pharmacy_id, product_id, unit_id);

CREATE INDEX product_unit_conversions_pharmacy_idx
  ON product_unit_conversions (pharmacy_id);

CREATE INDEX product_unit_conversions_pharmacy_and_product_idx
  ON product_unit_conversions (pharmacy_id, product_id);

-- Re-add the foreign key
ALTER TABLE product_unit_conversions
  ADD CONSTRAINT product_unit_conversions_pharmacy_id_fkey
  FOREIGN KEY (pharmacy_id) REFERENCES pharmacies(id);