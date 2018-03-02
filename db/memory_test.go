package db

import (
	"reflect"
	"testing"
)

func TestMemoryListEmpty(t *testing.T) {
	expectedKeys := []string{}
	db := NewMemoryDatabase()

	keys, err := db.List()
	if err != nil {
		t.Errorf("got error %q, want none", err)
	}

	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("got keys %q, want %q", keys, expectedKeys)
	}
}

func TestMemoryList(t *testing.T) {
	store := map[string]string{
		"key1": "value1",
		"key2": "value2",
	}
	expectedKeys := []string{"key1", "key2"}

	db := &memoryDatabase{
		store: store,
	}

	keys, err := db.List()
	if err != nil {
		t.Errorf("got error %q, want none", err)
	}

	if !reflect.DeepEqual(keys, expectedKeys) {
		t.Errorf("got keys %q, want %q", keys, expectedKeys)
	}
}

func TestMemoryGet(t *testing.T) {
	tests := []struct {
		desc  string
		store map[string]string
		key   string
		value string
		found bool
		err   error
	}{
		{
			desc: "success",
			store: map[string]string{
				"key1": "value1",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			db := &memoryDatabase{
				store: test.store,
			}

			value, found, err := db.Get(test.key)

			if err != test.err {
				t.Errorf("got error %q, want none", err)
			}

			if err != nil {
				return
			}

			if value != test.value {
				t.Errorf("got value %q, want %q", value, test.value)
			}

			if found != test.found {
				t.Errorf("got found %v, want %v", found, test.found)
			}
		})
	}
}

func TestMemoryPut(t *testing.T) {
	db := NewMemoryDatabase()

	_, found, err := db.Get("testkey")
	if err != nil {
		t.Errorf("got error %q, want none", err)
	}

	if found {
		t.Errorf("got found %v, want false", found)
	}

	if err := db.Put("testkey", "testvalue"); err != nil {
		t.Errorf("got error %q, want none", err)
	}

	content, found, err := db.Get("testkey")
	if err != nil {
		t.Errorf("got error %q, want none", err)
	}

	if !found {
		t.Errorf("got found %v, want true", found)
	}

	if content != "testvalue" {
		t.Errorf("got content %q, want %q", content, "testvalue")
	}
}
