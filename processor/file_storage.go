package processor

import (
	"github.com/bmatcuk/doublestar"
	"os"
)

//go:generate counterfeiter . StorageInterface
type StorageInterface interface {
	Write(path string, data []byte) error
	FileExists(path string) bool
	ExpandPaths(originalPaths []string) ([]string, error)
}

type FileStorage struct {
}

func (FileStorage) Write(path string, data []byte) error {
	out, _ := os.Create(path)
	defer out.Close()

	_, err := out.Write(data)
	return err
}

func (FileStorage) FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

func (FileStorage) ExpandPaths(originalPaths []string) ([]string, error) {
	resolvedPaths := []string{}

	for _, originalPath := range originalPaths {
		expanded, err := doublestar.Glob(originalPath)

		if err != nil {
			return []string{}, err
		} else {
			resolvedPaths = append(resolvedPaths, expanded...)
		}
	}

	return resolvedPaths, nil
}
