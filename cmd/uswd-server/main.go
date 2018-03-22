package main

import (
	"log"
	"net/http"

	"github.com/spf13/pflag"
	"github.com/xperimental/uswd/db"
	"github.com/xperimental/uswd/web"
)

var (
	baseDir = "./data/"
	addr    = ":8080"
)

func main() {
	pflag.StringVarP(&baseDir, "base", "b", baseDir, "Base directory of database.")
	pflag.StringVarP(&addr, "addr", "a", addr, "Network address to listen on.")
	pflag.Parse()

	db, err := db.NewFileDatabase(baseDir)
	if err != nil {
		log.Fatalf("Error initializing database: %s", err)
	}

	http.Handle("/", web.DatabaseHandler(db))

	log.Printf("Listening on %s...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
