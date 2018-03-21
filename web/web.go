package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/xperimental/uswd"
)

// NewRouter creates a new web router with all handlers.
func NewRouter(db uswd.Database) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGet(db, w, r)
		case http.MethodPut:
			handlePut(db, w, r)
		default:
			http.Error(w, fmt.Sprintf("Unknown method: %s", r.Method), http.StatusMethodNotAllowed)
		}
	})
}

func handleGet(db uswd.Database, w http.ResponseWriter, r *http.Request) {
	key := getKey(r)
	if key == "" {
		handleGetList(db, w, r)
		return
	}

	handleGetSingle(db, key, w, r)
}

func handleGetList(db uswd.Database, w http.ResponseWriter, r *http.Request) {
	keys, err := db.List()
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %s", err), http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(keys); err != nil {
		http.Error(w, fmt.Sprintf("Error encoding JSON: %s", err), http.StatusInternalServerError)
		return
	}
}

func handleGetSingle(db uswd.Database, key string, w http.ResponseWriter, r *http.Request) {
	content, err := db.Get(key)
	if _, ok := err.(uswd.NotFoundError); ok {
		http.Error(w, fmt.Sprintf("%s", err), http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Error getting content: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprint(w, content)
}

func handlePut(db uswd.Database, w http.ResponseWriter, r *http.Request) {
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

	if err := db.Put(key, string(content)); err != nil {
		http.Error(w, fmt.Sprintf("Error writing content: %s", err), http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, "saved.")
}

func getKey(r *http.Request) string {
	return strings.TrimPrefix(r.URL.Path, "/")
}
