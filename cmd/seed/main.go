package seed

import (
	"context"
	"log"

	"github.com/danivideda/satu-apotek-be/internal/config"
	"github.com/danivideda/satu-apotek-be/internal/db"
	"github.com/danivideda/satu-apotek-be/internal/dbsqlc"
)

func main() {
	cfg := config.Load()
	conn, err := db.NewPostgres(cfg.DB.URL)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()
	log.Println("Database connection pool established")

	ctx := context.Background()
	q := dbsqlc.New(conn)
	q.ProductUnitsCreate(ctx, dbsqlc.ProductUnitsCreateParams{})
}