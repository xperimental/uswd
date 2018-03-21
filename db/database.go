// Package db provides the database backend implementation.
package db

// Database is a simple key-value store.
type Database interface {
	List() ([]string, error)
	Get(key string) (content string, found bool, err error)
	Put(key, value string) error
}
