package db

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

type fileDatabase struct {
	baseDir string
}

// NewFileDatabase creates a database with a filesystem backend.
func NewFileDatabase(baseDir string) (Database, error) {
	stat, err := os.Stat(baseDir)
	switch {
	case os.IsNotExist(err):
		return nil, fmt.Errorf("directory does not exist: %s", baseDir)
	case err != nil:
		return nil, fmt.Errorf("error checking directory: %s", err)
	}

	if !stat.IsDir() {
		return nil, fmt.Errorf("not a directory: %s", baseDir)
	}

	return &fileDatabase{
		baseDir: baseDir,
	}, nil
}

func (d *fileDatabase) List() ([]string, error) {
	infos, err := ioutil.ReadDir(d.baseDir)
	if err != nil {
		return nil, err
	}

	keys := []string{}
	for _, i := range infos {
		keys = append(keys, i.Name())
	}

	return keys, nil
}

func (d *fileDatabase) Get(key string) (string, error) {
	path := filepath.Join(d.baseDir, key)

	content, err := ioutil.ReadFile(path)
	switch {
	case os.IsNotExist(err):
		return "", NotFoundError(key)
	case err != nil:
		return "", err
	}

	return string(content), nil
}

func (d *fileDatabase) Put(key, value string) error {
	path := filepath.Join(d.baseDir, key)
	return ioutil.WriteFile(path, []byte(value), 0666)
}
