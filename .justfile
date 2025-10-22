migration_path := "./cmd/migrate/migrations"

# Lists recipe
default:
    just --list --unsorted

# Create migration
migrate-create name:
    @migrate create -ext sql -dir cmd/migrate/migrations -seq {{name}}

# Run migration
migrate-up:
    @migrate -path={{migration_path}} -database=${MIGRATE_URL} up

# Reverse migration 'n' times
migrate-down n='':
    @migrate -path={{migration_path}} -database=${MIGRATE_URL} down {{n}}

# Check current migration version
migrate-version:
    @migrate -path={{migration_path}} -database=${MIGRATE_URL} version

# Force migration to specified version
migrate-force version:
    @migrate -path={{migration_path}} -database=${MIGRATE_URL} force {{version}}

# Migrate to specified version
migrate-goto version:
    @migrate -path={{migration_path}} -database=${MIGRATE_URL} goto {{version}}

# Seed the database
seed:
    @go run cmd/migrate/seed/*.go