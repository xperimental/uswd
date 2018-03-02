package web

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/xperimental/uswd"
)

// NewRouter creates a new web router with all handlers.
func NewRouter(db uswd.Database) http.Handler {
	r := mux.NewRouter()
	r.Path("/").Methods(http.MethodGet).HandlerFunc(listHandler(db))
	r.Path("/{key}").Methods(http.MethodGet).HandlerFunc(getHandler(db))
	r.Path("/{key}").Methods(http.MethodPut).HandlerFunc(putHandler(db))
	return r
}

func listHandler(db uswd.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
}

func getHandler(db uswd.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

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
}

func putHandler(db uswd.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		key := vars["key"]

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
}
