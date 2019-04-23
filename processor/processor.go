package processor

import (
	"github.com/remove-bg/go/client"
	"log"
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
		outputPath := DetermineOutputPath(inputPath, settings)

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
