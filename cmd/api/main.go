package main

import (
	"log"
	"os"
)

func main() {
	cfg := config{
		addr: os.Getenv("ADDR"),
	}

	app := &application{
		config: cfg,
	}

	mux := app.mount()
	log.Fatal(app.run(mux))
}
