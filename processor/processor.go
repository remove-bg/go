package processor

import (
	"github.com/remove-bg/go/client"
	"log"
	"path"
	"path/filepath"
	"strings"
)

type Processor struct {
	APIKey     string
	Client     client.ClientInterface
	FileWriter fileWriterInterface
}

type Settings struct {
	OutputDirectory string
}

func (p Processor) Process(inputPaths []string, settings Settings) {
	for _, inputPath := range inputPaths {
		outputPath := determineOutputPath(inputPath, settings.OutputDirectory)

		p.processFile(inputPath, outputPath)
	}
}

func (p Processor) processFile(inputPath string, outputPath string) {
	params := map[string]string{}

	processedBytes, err := p.Client.RemoveFromFile(inputPath, p.APIKey, params)

	if err != nil {
		log.Fatal(err)
		return
	}

	err = p.FileWriter.Write(outputPath, processedBytes)

	if err != nil {
		log.Fatal(err)
		return
	}
}

func determineOutputPath(inputPath string, outputDirectory string) string {
	_, fileName := filepath.Split(inputPath)
	withoutExtension := strings.TrimSuffix(fileName, path.Ext(fileName))
	return filepath.Join(outputDirectory, withoutExtension+".png")
}
