package processor

import (
	"github.com/remove-bg/go/client"
	"net/http"
)

type Processor struct {
	APIKey   string
	Client   client.ClientInterface
	Storage  StorageInterface
	Prompt   PromptInterface
	Notifier NotifierInterface
}

type Settings struct {
	OutputDirectory            string
	ReprocessExisting          bool
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
		Storage:  FileStorage{},
		Prompt:   Prompt{},
		Notifier: NewNotifier(),
	}
}

func (p Processor) Process(inputPaths []string, settings Settings) {
	confirmation := p.confirmLargeBatch(inputPaths, settings)
	if !confirmation {
		return
	}

	totalImages := len(inputPaths)

	for index, inputPath := range inputPaths {
		outputPath := DetermineOutputPath(inputPath, settings)
		skipImage := p.Storage.FileExists(outputPath) && !settings.ReprocessExisting

		if skipImage {
			p.Notifier.Skip(inputPath, outputPath, index+1, totalImages)
			return
		}

		err := p.processFile(inputPath, outputPath, settings.ImageSettings)

		if err == nil {
			p.Notifier.Success(inputPath, index+1, totalImages)
		} else {
			p.Notifier.Error(err, inputPath, index+1, totalImages)
		}
	}
}

func (p Processor) processFile(inputPath string, outputPath string, imageSettings ImageSettings) error {
	params := imageSettingsToParams(imageSettings)
	processedBytes, err := p.Client.RemoveFromFile(inputPath, p.APIKey, params)

	if err != nil {
		return err
	}

	return p.Storage.Write(outputPath, processedBytes)
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
	skipConfirm := settings.LargeBatchConfirmThreshold < 0

	if skipConfirm || batchSize < settings.LargeBatchConfirmThreshold {
		return true
	}

	return p.Prompt.ConfirmLargeBatch(batchSize)
}
