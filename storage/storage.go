package storage

import (
	"github.com/bmatcuk/doublestar"
	"os"
	"strings"
)

//go:generate counterfeiter . StorageInterface
type StorageInterface interface {
	Write(path string, data []byte) error
	FileExists(path string) bool
	ExpandPaths(originalPaths []string) ([]string, error)
	MkdirP(path string) error
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
		if !strings.Contains(originalPath, "*") {
			resolvedPaths = append(resolvedPaths, originalPath)
			continue
		}

		expanded, err := doublestar.Glob(originalPath)

		if err != nil {
			return []string{}, err
		} else {
			resolvedPaths = append(resolvedPaths, expanded...)
		}
	}

	return resolvedPaths, nil
}

func (FileStorage) MkdirP(path string) error {
	if len(path) == 0 {
		return nil
	}

	return os.MkdirAll(path, 0755)
}
