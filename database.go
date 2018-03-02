package uswd

// Database is a simple key-value store.
type Database interface {
	List() ([]string, error)
	Get(key string) (string, error)
	Put(key, value string) error
}
