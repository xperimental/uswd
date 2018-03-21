package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/xperimental/uswd/db"
)

// NewRouter creates a new web router with all handlers.
func NewRouter(database db.Database) http.Handler {
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
	key := getKey(r)
	if key == "" {
		handleGetList(database, w, r)
		return
	}

	handleGetSingle(database, key, w, r)
}

func handleGetList(database db.Database, w http.ResponseWriter, r *http.Request) {
	keys, err := database.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %s", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(keys); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %s", err), http.StatusInternalServerError)
		return
	}
}

func handleGetSingle(database db.Database, key string, w http.ResponseWriter, r *http.Request) {
	content, err := database.Get(key)
	if _, ok := err.(db.NotFoundError); ok {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting content: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, content)
}

func handlePut(database db.Database, w http.ResponseWriter, r *http.Request) {
	key := getKey(r)
	if key == "" {
		http.Error(w, "Key can not be empty!", http.StatusBadRequest)
		return
	}

	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error reading body: %s", err), http.StatusBadRequest)
		return
	}

	if err := database.Put(key, string(content)); err != nil {
		http.Error(w, fmt.Sprintf("Error writing content: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "saved.")
}

func getKey(r *http.Request) string {
	return strings.TrimPrefix(r.URL.Path, "/")
}
