package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var addr = ":8080"

func main() {
	flag.StringVar(&addr, "addr", addr, "Network address to listen on.")
	flag.Parse()

	http.HandleFunc("/", helloHandler)

	log.Printf("Listening on %s...", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello HTTP!")
}
