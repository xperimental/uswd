package web

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/xperimental/uswd/db"
)

type testDatabase struct {
	db  map[string]string
	err error
}

func (d *testDatabase) List() ([]string, error) {
	if d.err != nil {
		return nil, d.err
	}

	keys := []string{}
	for k := range d.db {
		keys = append(keys, k)
	}

	return keys, nil
}

func (d *testDatabase) Get(key string) (string, error) {
	if d.err != nil {
		return "", d.err
	}

	value, ok := d.db[key]
	if !ok {
		return "", db.NotFoundError(key)
	}

	return value, nil
}

func (d *testDatabase) Put(key, value string) error {
	if d.err != nil {
		return d.err
	}

	d.db[key] = value
	return nil
}

func TestRouterUnknownMethod(t *testing.T) {
	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodDelete, "/", nil)

	handler := DatabaseHandler(&testDatabase{})
	handler.ServeHTTP(w, r)

	expected := http.StatusMethodNotAllowed
	if w.Code != expected {
		t.Errorf("got status %d, want %d", w.Code, expected)
	}
}

func TestHandleGetList(t *testing.T) {
	for _, test := range []struct {
		desc string
		db   db.Database
		code int
		body string
	}{
		{
			desc: "empty",
			db: &testDatabase{
				db: map[string]string{},
			},
			code: http.StatusOK,
			body: "[]\n",
		},
		{
			desc: "one",
			db: &testDatabase{
				db: map[string]string{
					"key": "value",
				},
			},
			code: http.StatusOK,
			body: "[\"key\"]\n",
		},
		{
			desc: "error",
			db: &testDatabase{
				err: errors.New("test error"),
			},
			code: http.StatusInternalServerError,
			body: "Database error: test error\n",
		},
	} {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)

			handler := DatabaseHandler(test.db)
			handler.ServeHTTP(w, r)

			if w.Code != test.code {
				t.Errorf("got %d, want %d", w.Code, test.code)
			}

			if w.Body.String() != test.body {
				t.Errorf("got %q, want %q", w.Body.String(), test.body)
			}
		})
	}
}

func TestHandleGetSingle(t *testing.T) {
	for _, test := range []struct {
		desc string
		db   db.Database
		path string
		code int
		body string
	}{
		{
			desc: "success",
			db: &testDatabase{
				db: map[string]string{
					"key": "value",
				},
			},
			path: "/key",
			code: http.StatusOK,
			body: "value",
		},
		{
			desc: "not found",
			db: &testDatabase{
				db: map[string]string{},
			},
			path: "/key",
			code: http.StatusNotFound,
			body: "not found: key\n",
		},
		{
			desc: "error",
			db: &testDatabase{
				err: errors.New("test error"),
			},
			path: "/key",
			code: http.StatusInternalServerError,
			body: "Error getting content: test error\n",
		},
	} {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, test.path, nil)

			handler := DatabaseHandler(test.db)
			handler.ServeHTTP(w, r)

			if w.Code != test.code {
				t.Errorf("got status %d, want %d", w.Code, test.code)
			}

			if w.Body.String() != test.body {
				t.Errorf("got body %q, want %q", w.Body.String(), test.body)
			}
		})
	}
}

func TestHandlePut(t *testing.T) {
	for _, test := range []struct {
		desc  string
		db    db.Database
		path  string
		value string
		code  int
		body  string
	}{
		{
			desc: "success",
			db: &testDatabase{
				db: map[string]string{},
			},
			path:  "/key",
			value: "value",
			code:  http.StatusOK,
			body:  "saved.\n",
		},
		{
			desc: "no key",
			db: &testDatabase{
				db: map[string]string{},
			},
			path:  "/",
			value: "",
			code:  http.StatusBadRequest,
			body:  "Key can not be empty!\n",
		},
		{
			desc: "error",
			db: &testDatabase{
				err: errors.New("test error"),
			},
			path:  "/key",
			value: "",
			code:  http.StatusInternalServerError,
			body:  "Error writing content: test error\n",
		},
	} {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			w := httptest.NewRecorder()
			body := bytes.NewBufferString(test.value)
			r := httptest.NewRequest(http.MethodPut, test.path, body)

			handler := DatabaseHandler(test.db)
			handler.ServeHTTP(w, r)

			if w.Code != test.code {
				t.Errorf("got status %d, want %d", w.Code, test.code)
			}

			bodyStr := w.Body.String()
			if bodyStr != test.body {
				t.Errorf("got body %q, want %q", bodyStr, test.body)
			}
		})
	}
}
