package processor

import "os"

//go:generate counterfeiter . fileWriterInterface
type fileWriterInterface interface {
	Write(path string, data []byte) error
}

type FileWriter struct {
}

func (FileWriter) Write(path string, data []byte) error {
	out, _ := os.Create(path)
	defer out.Close()

	_, err := out.Write(data)
	return err
}
