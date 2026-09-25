# Layout and file names

```text
cmd/api/                         HTTP process. main.go wires, api.go mounts routes.
cmd/cron/                        gocron process. Off unless RUN_JOB=true.
cmd/migrate/migrations/          golang-migrate pairs. The only schema source.
internal/config/                 config.Load(). One file.
internal/env/                    getenv helpers. No business meaning.
internal/db/                     pgxpool open/close. Not sqlc.
internal/dbsqlc/queries/*.sql    queries you write.
internal/dbsqlc/*.go             sqlc output. Do not edit.
internal/repository/             sqlc wrappers, sentinels, session cache.
internal/service/                auth and session rules shared with middleware.
internal/http/handler/           HTTP adapters.
internal/http/middleware/        auth, CSRF, ownership guards.
internal/http/json/              read body, write {"data"} / {"error"}.
internal/http/cookie/            cookie names and flags.
internal/http/csrf/              HMAC token tied to a session id.
internal/http/jwt/               leftover. Do not extend.
internal/job/                    cron tasks. Called from cmd/cron only.
internal/tools/                  pinned CLI module (sqlc, migrate, air).
scripts/playground/              scratch. Not imported by the server.
```

`internal/db/db.go` opens the pool. `internal/dbsqlc/db.go` is generated sqlc. Do not confuse them.

## Naming a new aggregate

Use the plural snake name of the table. `products` is the example. Copy this set, do not invent synonyms (`product_controller.go`, `productRepo.go`).

| Piece | Path | Go name |
| --- | --- | --- |
| SQL | `internal/dbsqlc/queries/products.sql` | `-- name: GetProductByID :one` |
| Repository file | `internal/repository/products.go` | `type productsRepo struct`, interface `ProductsRepository` |
| Interface field | `Repository.Products` in `repository.go` | wired in `New` |
| Handler file | `internal/http/handler/products.go` | unexported `productHandler`, ctor `newProductHandler` |
| Handler field | `Handler.Product` in `handler.go` | methods exported: `Create`, `GetByID` |
| Service file | only if [layers.md](layers.md) says you need one | `internal/service/product.go` |
| Migration | `cmd/migrate/migrations/000010_<snake>.up.sql` and `.down.sql` | from `just migrate-create` |

Query names are `VerbNoun`, exported, with a sqlc annotation of `:one`, `:many`, or `:exec`. Parameters on new queries are `@snake_case`. Leave existing `$1` queries alone.

JSON fields and payload struct tags are `snake_case`. Add `validate` tags (`required`, `min`, `max`) on fields the handler must reject when empty. `parseAndValidateJSONPayload` already runs `go-playground/validator`.

One aggregate per file. Do not split a single handler into `products_create.go` unless the file is already large the way `pharmacies.go` is, and even then prefer staying in one file until it hurts.

## Names that already exist — do not "fix"

| Existing file | Leave it |
| --- | --- |
| `repository/owner_session.go` | singular, unlike the other session files |
| `repository/pharmacy_sessions.go`, `user_sessions.go` | plural |
| `handler/owners.go`, `service/owner.go` | plural vs singular across layers |
| `internal/http/jwt/*` | unused on purpose |

New session-like storage follows `pharmacy_sessions.go` (plural file, `PharmacySessionsRepository`).

## Route path shape

Prefix is always `/v1`. Group by actor, not by table name alone.

| Actor | Prefix | Example already in the tree |
| --- | --- | --- |
| Public / login | `/v1/auth/<actor>` | `/v1/auth/owners/login` |
| Owner | `/v1/owner/...` | `/v1/owner/pharmacies/{appID}/code/create` |
| Pharmacy device | `/v1/pharmacy/...` | `/v1/pharmacy/landing` |
| Staff user | `/v1/user/...` | `/v1/user/users/profile` |

Path segments are lowercase. A create is `/create`. A save is `/edit/save`. A single resource under an owner uses `{appID}` (Chi param), not `{id}`.
