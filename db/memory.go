package db

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
	// TODO implement list
	return []string{}, nil
}

func (db *memoryDatabase) Get(key string) (string, bool, error) {
	// TODO implement get
	return "", false, nil
}

func (db *memoryDatabase) Put(key, value string) error {
	// TODO implement put
	return nil
}
