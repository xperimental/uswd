package web

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/xperimental/uswd"
)

type testDatabase struct {
	db map[string]string
}

func (d *testDatabase) List() ([]string, error) {
	keys := []string{}
	for k := range d.db {
		keys = append(keys, k)
	}

	return keys, nil
}

func (d *testDatabase) Get(key string) (string, error) {
	value, ok := d.db[key]
	if !ok {
		return "", uswd.NotFoundError(key)
	}

	return value, nil
}

func (d *testDatabase) Put(key, value string) error {
	d.db[key] = value
	return nil
}

func TestListHandlerEmpty(t *testing.T) {
	for _, test := range []struct {
		desc string
		db   map[string]string
		code int
		body string
	}{
		{
			desc: "empty",
			db:   map[string]string{},
			code: http.StatusOK,
			body: "[]\n",
		},
		{
			desc: "one",
			db: map[string]string{
				"key": "value",
			},
			code: http.StatusOK,
			body: "[\"key\"]\n",
		},
	} {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			db := &testDatabase{
				db: test.db,
			}

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			handler := listHandler(db)
			handler(w, r)

			if w.Code != test.code {
				t.Errorf("got %d, want %d", w.Code, test.code)
			}

			if w.Body.String() != test.body {
				t.Errorf("got %q, want %q", w.Body.String(), test.body)
			}
		})
	}
}

func TestGetHandler(t *testing.T) {
	for _, test := range []struct {
		desc string
		db   map[string]string
		key  string
		code int
		body string
	}{
		{
			desc: "success",
			db: map[string]string{
				"key": "value",
			},
			key:  "key",
			code: http.StatusOK,
			body: "value",
		},
		{
			desc: "not found",
			db:   map[string]string{},
			key:  "key",
			code: http.StatusNotFound,
			body: "not found: key\n",
		},
	} {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			db := &testDatabase{
				db: test.db,
			}

			w := httptest.NewRecorder()
			r := mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/", nil), map[string]string{
				"key": test.key,
			})

			handler := getHandler(db)
			handler(w, r)

			if w.Code != test.code {
				t.Errorf("got status %d, want %d", w.Code, test.code)
			}

			if w.Body.String() != test.body {
				t.Errorf("got body %q, want %q", w.Body.String(), test.body)
			}
		})
	}
}

func TestPutHandler(t *testing.T) {
	for _, test := range []struct {
		desc  string
		db    map[string]string
		key   string
		value string
		code  int
	}{
		{
			desc:  "success",
			db:    map[string]string{},
			key:   "key",
			value: "value",
			code:  http.StatusOK,
		},
	} {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			db := &testDatabase{
				db: test.db,
			}

			w := httptest.NewRecorder()
			body := bytes.NewBufferString(test.value)
			r := mux.SetURLVars(httptest.NewRequest(http.MethodGet, "/", body), map[string]string{
				"key": test.key,
			})

			handler := putHandler(db)
			handler(w, r)

			if w.Code != test.code {
				t.Errorf("got status %d, want %d", w.Code, test.code)
			}
		})
	}
}
