CREATE TABLE IF NOT EXISTS sales (
  id bigserial PRIMARY KEY,
  pharmacy_id bigint NOT NULL REFERENCES pharmacies (id),

  total_price numeric(18,2) NOT NULL,
  is_pending boolean NOT NULL DEFAULT false,

  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),

  UNIQUE (pharmacy_id, id)
);
CREATE INDEX sales_pharmacy_id_created_at_idx ON sales (pharmacy_id, created_at DESC);

CREATE TABLE IF NOT EXISTS sale_items (
  id bigserial PRIMARY KEY,
  pharmacy_id bigint NOT NULL,

  sale_id bigint NOT NULL,
  product_id bigint REFERENCES products (id) ON DELETE SET NULL,
  
  -- some hard-copied value from product_id for records
  -- so if in the future some product is deleted, it's still recorded
  product_name text NOT NULL,
  qty_in_base_unit int NOT NULL,
  price_per_base_unit numeric(18,2) NOT NULL,

  selected_unit text NOT NULL,
  qty_in_selected_unit int NOT NULL,
  price_per_selected_unit numeric(18,2) NOT NULL,

  UNIQUE (sale_id, product_id),
  FOREIGN KEY (pharmacy_id, sale_id) REFERENCES sales (pharmacy_id, id) ON DELETE CASCADE
);
CREATE INDEX sale_items_pharmacy_id_product_id_idx ON sale_items (pharmacy_id, product_id);

