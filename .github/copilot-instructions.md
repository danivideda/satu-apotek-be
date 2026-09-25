Follow [AGENTS.md](../AGENTS.md) and the skill at
`.agents/skills/backend-structure/SKILL.md` before adding routes, handlers,
repositories, sqlc queries, or migrations.

Do not introduce another HTTP framework, ORM, or folder layout. This API is
Chi + pgx + sqlc, `GET`/`POST` only, three session actors (owner, pharmacy, user).
