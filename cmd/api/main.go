package main

import (
	"log"
	"os"

	"github.com/danivideda/satu-apotek-be/internal/db"
	"github.com/danivideda/satu-apotek-be/internal/http/handler"
)

func main() {
	cfg := config{
		addr: os.Getenv("ADDR"),
		db: dbConfig{
			url: os.Getenv("DATABASE_URL"),
		},
	}

	db, err := db.New(cfg.db.url)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Println("Database connection pool established")

	handler := handler.New(db)

	app := &application{
		config: cfg,
		handler: handler,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
