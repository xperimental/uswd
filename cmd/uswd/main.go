package main

import (
	"log"
	"net/http"

	"github.com/spf13/pflag"
	"github.com/xperimental/uswd"
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
