package processor

import "os"

//go:generate counterfeiter . StorageInterface
type StorageInterface interface {
	Write(path string, data []byte) error
	FileExists(path string) bool
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
