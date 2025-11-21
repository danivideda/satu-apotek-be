package main

import (
	"log"

	"github.com/danivideda/satu-apotek-be/internal/db"
	"github.com/danivideda/satu-apotek-be/internal/env"
	"github.com/danivideda/satu-apotek-be/internal/http/handler"
	"github.com/danivideda/satu-apotek-be/internal/repository"
)

func main() {
	cfg := config{
		addr: env.GetString("ADDR", "localhost:8080"),
		db: dbConfig{
			url: env.GetString("DATABASE_URL", "admin:adminpassword@localhost/satuapotek?sslmode=disable"),
		},
	}

	db, err := db.New(cfg.db.url)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Println("Database connection pool established")

	repository := repository.NewRepository(db)

	handler := handler.New(repository)

	app := &application{
		config: cfg,
		handler: handler,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
