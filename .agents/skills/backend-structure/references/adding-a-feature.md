# Adding a feature

Work in this order. Skip a step only when it already exists. The names below use `products` because that table is already migrated and not yet exposed. Swap the noun if the task is a different aggregate. Do not create the product files as a side effect of some other task.

1. **Schema.** If the table is missing, `just migrate-create <snake_name>`, write `.up.sql` and `.down.sql`, run `just migrate-up`. Products and sales already have migrations. Do not add a second `CREATE TABLE products`.
2. **Query.** Add `internal/dbsqlc/queries/products.sql`. Include `pharmacy_id` (or `owner_id`) in the SQL. Run `sqlc generate` and commit the generated Go.
3. **Repository.** Add `internal/repository/products.go`. Declare `ProductsRepository` and the `Products` field in `repository.go`. Construct it in `New`. Map pgx errors to `ErrNotFound` and `ErrDuplicateValue`.
4. **Service.** Skip unless middleware and the handler must share the rule. Products do not need one for a plain create/list.
5. **Handler.** Add `internal/http/handler/products.go` with an unexported struct and `newProductHandler`. Register the field in `Handler` and `New` in `handler.go`. Use `parseAndValidateJSONPayload`. Map sentinels to the `json.Response*` helpers.
6. **Route.** In `cmd/api/api.go` only, hang the method on the actor group that should call it.
   - Owner managing a shop: inside `/owner/pharmacies/{appID}`, so `GuardPharmacyDetailByOwner` has already run. Take the pharmacy from `PharmacyDetailFromCtx`.
   - Staff at the counter: inside `/user`, which already requires pharmacy session, user session, and CSRF.
   - Do not expose catalog writes on `/auth` or on a public route.
7. **Config.** Only if a new env var is required. Then `internal/env`, `config.Load`, and `.envrc.example` in the same change.
8. **Job.** Only for expiry or cleanup. Add `Add...Job` on `MyScheduler` and call it from `cmd/cron/main.go` next to the existing two jobs.

## Done check

- No new top-level directory and no new framework.
- No edits under `internal/dbsqlc` except generated output and `queries/*.sql`.
- No PUT, PATCH, or DELETE route.
- The handler never trusts a client-supplied `pharmacy_id` or `owner_id` for authorization.
- `gofmt` has been run on the Go files.
- A create returns `201` with `{"data": ...}`. A missing row returns `404` with `{"error":"not found"}`. A unique violation returns `400`, not `500`.

## Out of scope unless the task says so

Formatting the whole repo, renaming `owner_session.go`, deleting `internal/http/jwt`, filling the unimplemented `GetAllByAppID` stub, writing the missing `cmd/migrate/seed` program, or "fixing" `CSRFProtectionUser` (it currently clears the owner CSRF cookie).
