---
name: backend-structure
description: >
  Lay out and extend the Satu Apotek Go API (Chi, sqlc, pgx, session cookies).
  Use when adding or changing an endpoint, handler, service, repository, sqlc
  query, migration, middleware, route, cookie, or background job, or when
  deciding where a new file goes and what to name it.
---

# Backend structure

This repo is a small Chi service. Match the files that already exist. Do not
introduce a cleaner architecture beside them.

Module path: `github.com/danivideda/satu-apotek-be`.

## Read next

Open one of these only when the task needs it. Do not load them all up front.

| Task | File |
| --- | --- |
| Where a file lives, what to name it | [references/layout.md](references/layout.md) |
| Handler vs service vs repository | [references/layers.md](references/layers.md) |
| Routes, actors, cookies, JSON | [references/http.md](references/http.md) |
| Migrations, sqlc, tenant columns | [references/persistence.md](references/persistence.md) |
| Checklist for a new resource | [references/adding-a-feature.md](references/adding-a-feature.md) |

## Rules that override a "nicer" design

1. Routes are registered only in [`cmd/api/api.go`](../../../cmd/api/api.go). `main.go` only wires config, db, repository, service, handler, middleware.
2. HTTP methods are **GET and POST**. CORS allows those two. An update is a POST such as `/edit/save`. Do not add PUT, PATCH, or DELETE.
3. `/v1` requires `Content-Type: application/json`. Success is `{"data": ...}` from `json.WriteResponse`, `json.ResponseOK`, or `json.ResponseCreated`. Errors are `{"error": "<short generic>"}` from the `json.Response*` helpers. Do not add a second envelope, and do not return `err.Error()` to the client.
4. Three actors, three cookies. Do not merge them.
   - Owner: `owner_session` (HttpOnly) + `owner_csrf`
   - User (staff): `user_session` (HttpOnly) + `user_csrf`
   - Pharmacy (counter device): `pharmacy_session` (HttpOnly), no CSRF cookie
5. Mutating owner routes sit on the group that already applies `AuthOwner` and `CSRFProtectionOwner`. Mutating user routes sit on `AuthPharmacy` + `AuthUser` + `CSRFProtectionUser`. Do not remove CSRF to make a call succeed.
6. Pharmacy-owned rows are filtered with the pharmacy id from `middleware.AuthPharmacyFromCtx` or `middleware.PharmacyDetailFromCtx`. A pharmacy id in the JSON body is not authorization.
7. An owner addressing one pharmacy uses `/{appID}` plus `GuardPharmacyDetailByOwner`. `app_id` is the public id. Do not put the numeric pharmacy id in an owner URL.
8. Edit SQL in `internal/dbsqlc/queries/`, then run `sqlc generate`. Never hand-edit `internal/dbsqlc/*.go`.
9. Never edit a migration that has already been applied. Create the next one with `just migrate-create <snake_name>` and fill both `.up.sql` and `.down.sql`.
10. `internal/http/jwt` is unused. Do not mount it. Passwords use `argon2id` only.
11. New env vars go through `internal/env` into `config.Load()` and `.envrc.example`. Do not add a `.env` file.
12. A new cron job is a method on `job.MyScheduler`, registered in `cmd/cron/main.go`, gated by `RUN_JOB`. Do not start gocron from the HTTP process.
13. `scripts/playground` and `internal/tools` are not part of the API. Do not import them from `cmd/api`.
14. Do not rename or reformat unrelated files. Session filenames are already inconsistent; leave them.
15. Run `gofmt` on every Go file you touch.

## What not to add

No Gin, Echo, Fiber, GORM, ent, sqlx, or a `pkg/`, `domain/`, `usecase/`, `controller/`, or `entity/` tree. No new global middleware stack. No JWT access/refresh pair.
