# satu-apotek-be

Go API for Satu Apotek. Chi, PostgreSQL, sqlc, golang-migrate, session cookies.

Before adding or changing an endpoint, query, migration, or package, read
[`.agents/skills/backend-structure/SKILL.md`](.agents/skills/backend-structure/SKILL.md)
and only the reference it points at. Do not invent a second layout.

Hard stops, even if you have not opened the skill yet:

- No new router, ORM, or top-level folder (`pkg/`, `domain/`, `usecase/`).
- Routes are `GET` and `POST` only, under `/v1`, JSON in and `{"data":...}` / `{"error":...}` out.
- Three actors stay separate: owner, pharmacy, user. Do not merge their cookies or middleware.
- Do not edit generated files in `internal/dbsqlc/*.go` or already-applied migrations.
- Do not revive `internal/http/jwt`. Auth is session cookies + CSRF.
