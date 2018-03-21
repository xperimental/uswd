package db

import "fmt"

// Database is a simple key-value store.
type Database interface {
	List() ([]string, error)
	Get(key string) (string, error)
	Put(key, value string) error
}

// A NotFoundError is used when a key is not found in the database.
type NotFoundError string

func (e NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", string(e))
}
