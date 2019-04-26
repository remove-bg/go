package processor

import "os"

//go:generate counterfeiter . StorageInterface
type StorageInterface interface {
	Write(path string, data []byte) error
}

type FileStorage struct {
}

func (FileStorage) Write(path string, data []byte) error {
	out, _ := os.Create(path)
	defer out.Close()

	_, err := out.Write(data)
	return err
}
