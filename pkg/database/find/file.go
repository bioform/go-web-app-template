package find

import (
	"os"
	"path/filepath"
)

func File(path string) (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		absolutePath := filepath.Join(dir, path)
		if _, err := os.Stat(absolutePath); err == nil {
			return absolutePath, nil // found the root
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist // reached the filesystem root
		}
		dir = parent
	}
}
