# Persistence

PostgreSQL 17 via pgxpool. Queries are sqlc (`sqlc.yaml`: engine postgresql, `sql_package: pgx/v5`, `emit_json_tags: true`). Schema input for sqlc is `cmd/migrate/migrations`, not a live database.

## Migration

```bash
just migrate-create create_table_products
just migrate-up
just migrate-down 1
just migrate-version
```

`just migrate-create` writes the next sequential pair:

```text
cmd/migrate/migrations/000010_create_table_products.up.sql
cmd/migrate/migrations/000010_create_table_products.down.sql
```

The name is snake_case and says what changed (`create_table_products`, `alter_products_add_barcode`). Both files are required. The down script drops or reverts only what the up script added, in reverse order.

Do not change `000001` through the latest applied file. If a column is wrong, add a new migration.

Match the tables that are already there:

- Primary key `id bigserial`.
- Tenant column `pharmacy_id bigint NOT NULL REFERENCES pharmacies(id)` on anything that belongs to a shop. Index it (`<table>_pharmacy_idx`).
- Owner-level rows use `owner_id` and are read with the owner id from the session, the way pharmacies are.
- `created_at` and `updated_at` are `timestamptz NOT NULL DEFAULT NOW()`.
- Uniqueness is per pharmacy (`UNIQUE (pharmacy_id, name)`), not global, unless the column is a true global key (owner email, owner username).

`products`, `product_categories`, `product_units`, `product_unit_conversions` (migration `000008`), plus `sales` and `sale_items` (migration `000009`), already exist and have **no** sqlc queries yet. Wire those tables; do not recreate them.

## sqlc

One file per table family: `internal/dbsqlc/queries/products.sql`.

```sql
-- name: GetProductByIDForPharmacy :one
SELECT * FROM products
WHERE id = @id AND pharmacy_id = @pharmacy_id;

-- name: CreateProduct :one
INSERT INTO products (pharmacy_id, name)
VALUES (@pharmacy_id, @name)
RETURNING *;
```

Every read and write of a pharmacy-owned table includes `pharmacy_id` in the SQL. Do not fetch by `id` alone and filter in Go.

Then:

```bash
sqlc generate
```

Commit the generated `internal/dbsqlc/products.sql.go` (and any `models.go` diff) together with the `.sql`. If you did not run sqlc, do not hand-write the Go.

After generate:

1. Add methods on a new `ProductsRepository` interface in `repository.go`.
2. Implement them on `productsRepo` by calling `r.queries.GetProductByIDForPharmacy`.
3. Return `ErrNotFound` / `ErrDuplicateValue` at this boundary.

`SELECT *` is what current queries use. Keep that until the repo decides otherwise. Do not add a query the handler does not call.

## IDs

Numeric `id` stays inside the process. The value the owner frontend already uses for a pharmacy is `app_id` (sqids, written when the pharmacy is created). Do not accept `pharmacy_id` from an owner payload when `app_id` is the identifier that route family uses. Staff and pharmacy-device routes may use the internal id taken from the session, not from the client.
