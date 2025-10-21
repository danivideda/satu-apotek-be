package tools

import (
	_ "github.com/air-verse/air"                             // Live reload Go on change
	_ "github.com/golang-migrate/migrate/v4/cmd/migrate"     // create migration for SQL
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"                    // SQL safe type
)
