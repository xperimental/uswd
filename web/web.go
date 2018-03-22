// Package web provides HTTP handlers for the server.
package web

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/xperimental/uswd/db"
)

// DatabaseHandler creates a HTTP handler for interacting with a database.
func DatabaseHandler(database db.Database) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGet(database, w, r)
		case http.MethodPut:
			handlePut(database, w, r)
		default:
			http.Error(w, fmt.Sprintf("Unknown method: %s", r.Method), http.StatusMethodNotAllowed)
		}
	})
}

func handleGet(database db.Database, w http.ResponseWriter, r *http.Request) {
	// TODO implement get
}

func handleGetList(database db.Database, w http.ResponseWriter, r *http.Request) {
	// TODO implement get list
}

func handleGetSingle(database db.Database, key string, w http.ResponseWriter, r *http.Request) {
	// TODO implement get single
}

func handlePut(database db.Database, w http.ResponseWriter, r *http.Request) {
	// TODO implement put
}

func getKey(r *http.Request) string {
	return strings.TrimPrefix(r.URL.Path, "/")
}
