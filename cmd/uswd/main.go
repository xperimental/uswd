package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/xperimental/uswd"
	"github.com/xperimental/uswd/web"
)

var (
	baseDir = "./data/"
	addr    = ":8080"
)

func main() {
	flag.StringVar(&baseDir, "base", baseDir, "Base directory of database.")
	flag.StringVar(&addr, "addr", addr, "Network address to listen on.")
	flag.Parse()

	db, err := uswd.NewFileDatabase(baseDir)
	if err != nil {
		log.Fatalf("Error initializing database: %s", err)
	}

	router := web.NewRouter(db)

	log.Printf("Listening on %s...", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatal(err)
	}
}
