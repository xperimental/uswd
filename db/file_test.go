package db

import (
	"errors"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestNewFileDatabase(t *testing.T) {
	tests := []struct {
		desc string
		path string
		err  error
	}{
		{
			desc: "success",
			path: "_testdata",
			err:  nil,
		},
		{
			desc: "file",
			path: "file_test.go",
			err:  errors.New("not a directory: file_test.go"),
		},
		{
			desc: "not existing",
			path: "does-not-exist",
			err:  errors.New("directory does not exist: does-not-exist"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			_, err := NewFileDatabase(test.path)

			if !reflect.DeepEqual(err, test.err) {
				t.Errorf("got error %q, wanted %q", err, test.err)
			}
		})
	}
}

func TestFileList(t *testing.T) {
	expectedKeys := []string{"key1"}

	db, err := NewFileDatabase("_testdata")
	if err != nil {
		t.Fatalf("error creating database: %s", err)
	}

	keys, err := db.List()
	if err != nil {
		t.Fatalf("got error %q, wanted none", err)
	}

	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("got keys %q, wanted %q", keys, expectedKeys)
	}
}

func TestFileGet(t *testing.T) {
	tests := []struct {
		desc    string
		key     string
		content string
		found   bool
		err     error
	}{
		{
			desc:    "not found",
			key:     "not-found",
			content: "",
			found:   false,
			err:     nil,
		},
		{
			desc:    "success",
			key:     "key1",
			content: "value1",
			found:   true,
			err:     nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			db, err := NewFileDatabase("_testdata")
			if err != nil {
				t.Fatalf("error creating database: %s", err)
			}

			content, found, err := db.Get(test.key)

			if !reflect.DeepEqual(err, test.err) {
				t.Errorf("got error %q, wanted %q", err, test.err)
			}

			if err != nil {
				return
			}

			if content != test.content {
				t.Errorf("got content %q, wanted %q", content, test.content)
			}

			if found != test.found {
				t.Errorf("got found %v, want %v", found, test.found)
			}
		})
	}
}

func TestFilePut(t *testing.T) {
	tests := []struct {
		desc    string
		key     string
		content string
		found   bool
		err     error
	}{
		{
			desc:    "success",
			key:     "key1",
			content: "value1",
			found:   true,
			err:     nil,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			dir, err := ioutil.TempDir("", "uswd")
			if err != nil {
				t.Fatalf("error creating temporary directory: %s", err)
			}
			defer os.RemoveAll(dir)

			db, err := NewFileDatabase(dir)
			if err != nil {
				t.Fatalf("error creating database: %s", err)
			}

			err = db.Put(test.key, test.content)

			if !reflect.DeepEqual(err, test.err) {
				t.Errorf("got error %q, wanted %q", err, test.err)
			}
		})
	}
}
