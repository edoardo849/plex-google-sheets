package cache

import (
	"os"
	"os/user"
	"path/filepath"
)

// FilePath generates credential file path/filename.
// It returns the generated credential path/filename.
func FilePath(path, filename string) (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	// path is .credentials
	//
	cacheDir := filepath.Join(usr.HomeDir, path)
	os.MkdirAll(cacheDir, 0700)
	return filepath.Join(cacheDir, filename), err
}

// Open opens the cache file
func Open(file string) (*os.File, error) {

	return os.Open(file)
}

// OpenOrCreate opens or create a cache file
func OpenOrCreate(file string) (*os.File, error) {

	return os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
}
