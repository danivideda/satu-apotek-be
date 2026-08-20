CREATE TABLE IF NOT EXISTS product_categories (
  id bigserial PRIMARY KEY,
  pharmacy_id bigint NOT NULL REFERENCES pharmacies(id),

  code VARCHAR(10) NOT NULL,
  name text NOT NULL,

  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),

  UNIQUE (pharmacy_id, name)
);
CREATE INDEX product_categories_pharmacy_idx ON product_categories (pharmacy_id);

CREATE TABLE IF NOT EXISTS product_units (
  id bigserial PRIMARY KEY,
  pharmacy_id bigint NOT NULL REFERENCES pharmacies(id),

  label VARCHAR(10) NOT NULL,
  name text NOT NULL,

  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),

  UNIQUE (pharmacy_id, label),
  UNIQUE (pharmacy_id, name)
);
CREATE INDEX product_units_pharmacy_idx ON product_units (pharmacy_id);

CREATE TABLE IF NOT EXISTS products (
  id bigserial PRIMARY KEY,
  pharmacy_id bigint NOT NULL REFERENCES pharmacies(id),

  name text NOT NULL,
  sku text,
  barcode text,
  category_id bigint REFERENCES product_categories(id),
  base_unit_id bigint REFERENCES product_units(id),
  min_stock int NOT NULL DEFAULT 0,
  base_sell_price decimal(18,2) NOT NULL DEFAULT 0,
  is_active boolean NOT NULL DEFAULT true,

  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW()
);
CREATE INDEX products_pharmachy_idx ON products (pharmacy_id);
CREATE INDEX products_pharmachy_and_sku_idx ON products (pharmacy_id, sku);
CREATE INDEX products_pharmachy_and_barcode_idx ON products (pharmacy_id, barcode);
CREATE INDEX products_pharmachy_and_name_idx ON products (pharmacy_id, name);
CREATE INDEX products_pharmachy_and_category_idx ON products (pharmacy_id, category_id);

CREATE TABLE IF NOT EXISTS product_unit_conversions (
  id bigserial PRIMARY KEY,
  pharmacy_id bigint NOT NULL REFERENCES pharmacies(id),

  product_id bigint REFERENCES products(id),
  unit_id bigint REFERENCES product_units(id),
  quantity_in_base_unit integer NOT NULL DEFAULT 1,

  created_at timestamptz NOT NULL DEFAULT NOW(),
  updated_at timestamptz NOT NULL DEFAULT NOW(),

  UNIQUE (pharmacy_id, product_id, unit_id)
);
CREATE INDEX product_unit_conversions_pharmacy_idx ON product_unit_conversions (pharmacy_id);
CREATE INDEX product_unit_conversions_pharmacy_and_product_idx ON product_unit_conversions (pharmacy_id, product_id);
