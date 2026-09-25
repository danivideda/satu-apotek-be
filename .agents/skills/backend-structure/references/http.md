# HTTP

Router: Chi v5, mounted in `application.mount` in `cmd/api/api.go`.

Global middleware, in order: CORS, RequestID, RealIP, Logger, Recoverer, 60s timeout. Then `/v1` applies `AllowContentType("application/json")`.

CORS origins come from `CORS_ALLOW_ADDR` (default `http://localhost:3000` and `http://localhost:4173`). Methods are `GET` and `POST`. Headers are `X-CSRF-Token` and `Content-Type`. Credentials are allowed. Changing this list is a config change, not a per-route hack.

## Actors

| Actor | Who | Cookie | CSRF cookie | Middleware | Context helper |
| --- | --- | --- | --- | --- | --- |
| Owner | account that owns pharmacies | `owner_session` | `owner_csrf` | `AuthOwner`, `CSRFProtectionOwner` | `AuthOwnerFromCtx` |
| Pharmacy | the shop device, after a connect code | `pharmacy_session` | none | `AuthPharmacy` | `AuthPharmacyFromCtx` |
| User | staff inside one pharmacy | `user_session` | `user_csrf` | `AuthPharmacy` then `AuthUser`, `CSRFProtectionUser` | `AuthUserFromCtx` |

A staff request is pharmacy session **and** user session. Do not authorize a user route with only `AuthUser`.

`GuardPharmacyDetailByOwner` must wrap `/{appID}` owner routes. Handlers under that mount read `PharmacyDetailFromCtx` instead of loading the pharmacy again by an id from the body.

Login and register stay under `/v1/auth/...` and do not use the authenticated groups. Logout is authenticated and CSRF-protected. `GET .../check` is authenticated and not CSRF-protected, matching the current auth routes.

Pharmacy connect (`POST /v1/auth/pharmacies/connect`) is public on purpose: there is no session yet. Do not hang CSRF on it.

## Registering a route

Add it inside the existing `r.Route` group. Do not create a second `/v1` mount.

```go
r.Route("/owner", func(r chi.Router) {
    r.Use(app.middleware.AuthOwner)
    r.Use(app.middleware.CSRFProtectionOwner)

    r.Route("/products", func(r chi.Router) {
        r.Post("/create", app.handler.Product.Create)
    })
})
```

Chi params: `chi.URLParam(r, "appID")`. The param name matches the pattern, including casing.

## JSON

- Read with `parseAndValidateJSONPayload`. Unknown JSON fields fail (`DisallowUnknownFields`). Body cap is 1 MiB. Empty body is an error.
- Write success with `json.ResponseOK`, `json.ResponseCreated`, or `json.WriteResponse(w, status, payload)`. All three wrap the value as `{"data": ...}`. Do not call `json.Write` from a handler or the client loses the envelope.
- `ResponseNoContent` is available. Prefer `ResponseCreated` for a create, matching `users.Create`.
- Error text is fixed inside `internal/http/json/response.go` (`"bad request"`, `"not found"`, `"unauthorized"`, `"forbidden"`, `"invalid CSRF token"`, `"the server encountered some problem"`). Log the real error. Do not add the internal cause to the response.
- Health `GET /v1/health` is the one plain-text exception. Leave it.

## CSRF

Owner and user mutating routes inherit CSRF from the group `Use`. The middleware skips `GET`. `POST` must send the non-HttpOnly CSRF cookie value in `X-CSRF-Token`. Do not mark the CSRF cookie HttpOnly. Do not skip the middleware for a single POST inside that group. Keep new owner routes inside `/owner` and new staff routes inside `/user` so they pick up the same middleware.

On failure the middleware deletes the CSRF cookie so the next `/check` can refresh it. Do not copy the user middleware's current call to `cookie.Owner.DeleteCSRF` into new code, and do not change that line unless the task is that bug.
