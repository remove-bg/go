package processor

import (
	"github.com/remove-bg/go/client"
	"log"
	"net/http"
)

type Processor struct {
	APIKey     string
	Client     client.ClientInterface
	FileWriter fileWriterInterface
	Prompt     promptInterface
}

type Settings struct {
	OutputDirectory            string
	SkipConfirmLargeBatch      bool
	LargeBatchConfirmThreshold int
	ImageSettings              ImageSettings
}

type ImageSettings struct {
	Size     string
	Type     string
	Channels string
	BgColor  string
}

func NewProcessor(apiKey string) Processor {
	return Processor{
		APIKey: apiKey,
		Client: client.Client{
			HTTPClient: http.Client{},
		},
		FileWriter: FileWriter{},
		Prompt:     Prompt{},
	}
}

func (p Processor) Process(inputPaths []string, settings Settings) {
	confirmation := p.confirmLargeBatch(inputPaths, settings)
	if !confirmation {
		return
	}

	for _, inputPath := range inputPaths {
		outputPath := DetermineOutputPath(inputPath, settings)

		p.processFile(inputPath, outputPath, settings.ImageSettings)
	}
}

func (p Processor) processFile(inputPath string, outputPath string, imageSettings ImageSettings) {
	params := imageSettingsToParams(imageSettings)
	processedBytes, err := p.Client.RemoveFromFile(inputPath, p.APIKey, params)

	if err != nil {
		log.Print(err)
		return
	}

	err = p.FileWriter.Write(outputPath, processedBytes)

	if err != nil {
		log.Print(err)
		return
	}
}

func imageSettingsToParams(imageSettings ImageSettings) map[string]string {
	// TODO: Tidyup with reflection / struct tags?
	params := map[string]string{}

	if len(imageSettings.Size) > 0 {
		params["size"] = imageSettings.Size
	}

	if len(imageSettings.Type) > 0 {
		params["type"] = imageSettings.Type
	}

	if len(imageSettings.Channels) > 0 {
		params["channels"] = imageSettings.Channels
	}

	if len(imageSettings.BgColor) > 0 {
		params["bg_color"] = imageSettings.BgColor
	}

	return params
}

func (p Processor) confirmLargeBatch(inputPaths []string, settings Settings) bool {
	batchSize := len(inputPaths)
	overThreshold := batchSize < settings.LargeBatchConfirmThreshold

	if overThreshold || settings.SkipConfirmLargeBatch {
		return true
	}

	return p.Prompt.ConfirmLargeBatch(batchSize)
}
