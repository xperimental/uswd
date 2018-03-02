package uswd

import "fmt"

// Database is a simple key-value store.
type Database interface {
	List() ([]string, error)
	Get(key string) (string, error)
	Put(key, value string) error
}

type NotFoundError string

func (e NotFoundError) Error() string {
	return fmt.Sprintf("not found: %s", string(e))
}
