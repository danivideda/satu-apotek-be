# satu-apotek-be

Backend for **Satu Apotek**

Built with **Go + Chi**, PostgreSQL, sqlc, and golang-migrate.

## Tech Stack

| Layer | Choice |
|-------|--------|
| Language | Go 1.25.3 |
| HTTP Router | [Chi v5](https://github.com/go-chi/chi) |
| Database | PostgreSQL 17 |
| Migrations | [golang-migrate](https://github.com/golang-migrate/migrate) |
| Query generation | [sqlc](https://sqlc.dev) |
| Auth | Session cookies + CSRF (HMAC) + argon2id |
| Session cache | [go-cache](https://github.com/patrickmn/go-cache) |
| Live reload | [Air](https://github.com/air-verse/air) |
| Task runner | [Just](https://github.com/casey/just) |
| Env loading | [direnv](https://direnv.net/docs/installation.html) |

## Prerequisites

- **Go** 1.25.3+
- **Docker** + **Docker Compose**
- **direnv** ≥ 2.35 
- **Just** ≥ 1.41

## Quick Start

```bash
# 1. Clone
git clone https://github.com/danivideda/satu-apotek-be.git
cd satu-apotek-be

# 2. Environment
cp .envrc.example .envrc
direnv allow .

# 3. Install project tools (migrate, sqlc, air) into ./bin
chmod +x bin/install_tools
./bin/install_tools

# 4. Start Postgres
docker compose up -d

# 5. Run migrations
just migrate-up

# 6. Start the API (live reload)
air
# or without air:
# go run ./cmd/api/*.go
```

The server listens on `http://localhost:8080` by default. Health check: `GET /v1/health`

## Environment Variables

All variables are loaded via direnv from `.envrc`.  
See `.envrc.example` for the full list.

> [!NOTE]
> All secret values shown below are for **local development only**.

Key ones:

| Variable | Purpose | Default (local) |
|----------|---------|-----------------|
| `ADDR` | HTTP listen address | `:8080` |
| `DATABASE_URL` | Postgres connection string (pgx) | `postgres://admin:adminpassword@localhost/satuapotek?sslmode=disable` |
| `MIGRATE_URL` | Same DB but with `pgx5://` scheme for migrate | `pgx5://admin:...` |
| `OWNER_SESSION_TTL` | Owner session lifetime | `10m` |
| `USER_SESSION_TTL` | User session lifetime | `10m` |
| `PHARMACY_SESSION_TTL` | Pharmacy session lifetime | `10m` |
| `CACHE_SESSION_TTL` | In-memory cache TTL for sessions | `1m` |
| `CODE_TTL` | Pharmacy code lifetime | `1m` |
| `CSRF_SECRET` | Secret used to sign CSRF tokens | (see `.envrc.example`) |
| `RUN_JOB` | Enable background jobs | `false` |
| `CRON_DURATION_DEL_EXP_SESSION` | How often to clean expired sessions | `10s` |
| `CRON_DURATION_CLEAR_APTK_CODE` | How often to clear expired pharmacy codes | `10s` |

> [!NOTE]
> The JWT-related variables (`JWT_ACCESS_SECRET`, `JWT_REFRESH_SECRET`, `ACCESS_TTL`, `REFRESH_TTL`) still exist in `.envrc.example` but are currently unused. Auth is fully session + cookie based.

## Database

### Local Postgres (Docker Compose)

```bash
docker compose up -d          # starts postgres:17.6-alpine on :5432
docker compose down           # stop
docker compose down -v        # stop + wipe data volume
```

Credentials (from `compose.yaml`):

- DB: `satuapotek`
- User: `admin`
- Password: `adminpassword`

### Migrations

Migrations live in `cmd/migrate/migrations/` and are managed with the `migrate` CLI (installed by `./bin/install_tools`).

```bash
# Apply all pending migrations
just migrate-up

# Roll back the last migration
just migrate-down 1

# Check current version
just migrate-version

# Create a new sequential migration
just migrate-create add_something_table

# Force version (use with care)
just migrate-force 7
```

`just` recipes are defined in `.justfile`.

### sqlc

SQL queries live in `internal/dbsqlc/queries/`.  
Schema is taken from the migration files.

```bash
# Regenerate Go code after changing queries or schema
sqlc generate
```

Config: `sqlc.yaml`.

## Development Workflow

1. Make schema changes → create migration with `just migrate-create ...`
2. Write / adjust SQL in `internal/dbsqlc/queries/`
3. `sqlc generate`
4. Implement repository / service / handler layers
5. `air` for hot reload while developing

### Useful commands

```bash
just                  # list all recipes
just migrate-up
just migrate-down 1
just migrate-version

# Install / reinstall tools
./bin/install_tools
```

## Project Structure (high level)

```markdown
cmd/
├── api/                 # HTTP server entrypoint
└── migrate/             # migration files + (future) seed
internal/
├── db/                  # pgx connection helper
├── dbsqlc/              # generated sqlc code + raw queries
├── env/                 # env helpers
├── http/
│   ├── handler/         # HTTP handlers
│   ├── middleware/      # auth, CSRF, guards
│   ├── jwt/             # leftover JWT helpers (currently unused)
│   └── json/            # response helpers
├── repository/          # data access + in-memory session cache
├── service/             # business logic (cookies, CSRF, etc.)
└── job/                 # background jobs (gocron)
bin/                     # local tool binaries
```

## Notes

- Auth uses **HttpOnly session cookies** + **CSRF cookies**. Sessions are stored in Postgres and cached in-memory.
- Background jobs (expired session cleanup, apotek code expiry) are controlled by `RUN_JOB=true` and the `CRON_DURATION_*` variables.
- The `seed` recipe exists in `.justfile` but the seed implementation is not present yet.
- Frontend CORS is currently allowed for `http://localhost:3000` and `http://localhost:4173`.