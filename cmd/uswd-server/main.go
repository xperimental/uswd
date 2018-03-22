package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/xperimental/uswd/db"
	"github.com/xperimental/uswd/web"
)

var addr = ":8080"

func main() {
	flag.StringVar(&addr, "addr", addr, "Network address to listen on.")
	flag.Parse()

	db := db.NewMemoryDatabase()

	http.Handle("/", web.DatabaseHandler(db))

	log.Printf("Listening on %s...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
