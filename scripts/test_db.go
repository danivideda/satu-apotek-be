package main

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/danivideda/satu-apotek-be/internal/db"
)

func main() {
	dbUrl := os.Getenv("DATABASE_URL")
	parsedUrl, err := url.Parse(dbUrl)
	if err != nil {
		panic(err)
	}
	parsedUrl.Scheme = "postgres"
	dbUrl = parsedUrl.String()

	db, err := db.New(dbUrl)
	if err != nil {
		panic(err)
	}

	var greeting string
	err = db.QueryRow(context.Background(), `select 'hello world'`).Scan(&greeting)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Message: %s", greeting)
}
