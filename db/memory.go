package db

import "sort"

type memoryDatabase struct {
	store map[string]string
}

// NewMemoryDatabase creates a simple in-memory key-value store.
func NewMemoryDatabase() Database {
	return &memoryDatabase{
		store: make(map[string]string),
	}
}

func (db *memoryDatabase) List() ([]string, error) {
	keys := []string{}
	for k := range db.store {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	return keys, nil
}

func (db *memoryDatabase) Get(key string) (string, bool, error) {
	value, ok := db.store[key]
	return value, ok, nil
}

func (db *memoryDatabase) Put(key, value string) error {
	db.store[key] = value
	return nil
}
