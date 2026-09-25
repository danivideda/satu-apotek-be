# Layers

```text
chi route (cmd/api/api.go)
  -> middleware (session, CSRF, ownership guard)
    -> handler (JSON, validate, status)
      -> repository            for one query or one transaction
      -> service -> repository  when the rule is shared with middleware
        -> dbsqlc
          -> Postgres
```

Handlers call the repository directly today (`users.go`, `owners.go`, most of `pharmacies.go`). That is the default. Add a service only when the same rule must run from both a handler and middleware, or when the handler would otherwise own session rotation, cache fill, or a credential check. `service.Session` and `service.Auth` are the pattern. A service method that only forwards one repository call is noise. Do not add it.

## Handler

File: `internal/http/handler/<plural>.go`.

- Unexported struct (`userHandler`) holding `repository.Repository` and, only if needed, a `*service.Auth` or a TTL from `config.Config`.
- Constructor `newUserHandler`, called from `New` in `handler.go`.
- Exported methods with signature `func (h *userHandler) Create(w http.ResponseWriter, r *http.Request)`.
- Read the principal with `middleware.AuthOwnerFromCtx`, `AuthUserFromCtx`, `AuthPharmacyFromCtx`, or `PharmacyDetailFromCtx`. If the type assertion fails, that is a server bug: `json.ResponseInternalServerError`.
- Body: local payload struct, then `parseAndValidateJSONPayload`. It already answers 400.
- Map repository and service sentinels with `errors.Is`. Unknown errors are 500.

| Sentinel | Status helper |
| --- | --- |
| `repository.ErrNotFound` | `json.ResponseNotFound` |
| `repository.ErrDuplicateValue` | `json.ResponseBadRequest` |
| `service.ErrInvalidCredentials`, `service.ErrInvalidSession` | `json.ResponseUnauthorized` |
| `service.ErrUserForbidden` | `json.ResponseForbidden` |
| anything else | `json.ResponseInternalServerError` |

Do not compare `err.Error()` strings. Add a new sentinel next to the existing ones in `repository/errors.go` or `service/errors.go` when a new case needs its own status.

Cookie writes happen in the handler (login / logout) through the package-level `cookie` in `handler.go`, or in middleware when a session is rotated. Do not `http.SetCookie` with a new name.

## Repository

File: `internal/repository/<plural>.go`.

- Interface on `Repository` in `repository.go`: `Products ProductsRepository`.
- Concrete type unexported, constructed in `New`.
- Hold `*dbsqlc.Queries`. Hold `*pgxpool.Pool` as well only when the method needs a transaction (see `pharmaciesRepo`, which writes the row and then the sqids `app_id`).
- Translate `pgx.ErrNoRows` with `isNotFoundError` to `ErrNotFound`. Translate Postgres `23505` with `isDuplicateError` to `ErrDuplicateValue`. Do not leak `pgconn.PgError` to handlers.
- A transaction stays inside the repository method. Handlers do not begin transactions.
- `CacheStore` is for sessions. Do not cache pharmacy catalog rows there unless the task is explicitly about that cache.

## Service

File: `internal/service/<topic>.go`.

- Construct with `NewXxx(repo, cfg)` from `cmd/api/main.go`. Middleware builds its own `service.Session` inside `middleware.New`; keep that wiring, do not also stash a second session service on the handler unless `main` already passed one in (auth does).
- Return sentinels from `service/errors.go`, not HTTP statuses. The service package does not import `net/http`.

## Middleware

File stays `internal/http/middleware/`. New auth is not a new framework. Add a guard next to `guards.go` when a route must prove ownership before the handler runs, the way `GuardPharmacyDetailByOwner` loads the pharmacy and stores it on the context.

Context keys are unexported constants. Export a `SomethingFromCtx` function. Do not read `r.Context().Value("...")` from a handler with a string literal.

## Config

`config.Load()` is the only place that reads the environment. Handlers do not call `os.Getenv`. A duration or secret the handler needs is copied onto the handler struct in `New` (see `PharmacySessionTTL`).
