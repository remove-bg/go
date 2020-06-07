package processor

import (
	"fmt"
	"github.com/remove-bg/go/client"
	"github.com/remove-bg/go/composite"
	"github.com/remove-bg/go/storage"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type Processor struct {
	APIKey     string
	Client     client.ClientInterface
	Storage    storage.StorageInterface
	Prompt     PromptInterface
	Notifier   NotifierInterface
	Compositor composite.CompositorInterface
}

type Settings struct {
	OutputDirectory            string
	ReprocessExisting          bool
	SkipPngFormatOptimization  bool
	LargeBatchConfirmThreshold int
	ImageSettings              ImageSettings
}

type ImageSettings struct {
	Size            string
	Type            string
	Channels        string
	BgColor         string
	BgImageFile     string
	Format          string
	ExtraApiOptions string
}

func NewProcessor(apiKey string, version string) Processor {
	return Processor{
		APIKey: apiKey,
		Client: client.Client{
			Version:    version,
			HTTPClient: http.Client{},
		},
		Storage:    storage.FileStorage{},
		Prompt:     Prompt{},
		Notifier:   NewNotifier(),
		Compositor: composite.New(),
	}
}

func (p Processor) Process(rawInputPaths []string, settings Settings) {
	err := p.Storage.MkdirP(settings.OutputDirectory)
	if err != nil {
		log.Fatal(err)
	}

	inputPaths, err := p.Storage.ExpandPaths(rawInputPaths)
	if err != nil {
		log.Fatal(err)
	}

	confirmation := p.confirmLargeBatch(inputPaths, settings)
	if !confirmation {
		return
	}

	if !settings.SkipPngFormatOptimization {
		settings.ImageSettings.upgradePngToZipFormat()
	}

	totalImages := len(inputPaths)

	for index, inputPath := range inputPaths {
		outputPath := DetermineOutputPath(inputPath, settings)
		skipImage := p.Storage.FileExists(outputPath) && !settings.ReprocessExisting

		if skipImage {
			p.Notifier.Skip(inputPath, outputPath, index+1, totalImages)
			continue
		}

		err := p.processFile(inputPath, outputPath, settings.ImageSettings)

		if err == nil {
			p.Notifier.Success(inputPath, index+1, totalImages)
		} else {
			p.Notifier.Error(err, inputPath, index+1, totalImages)

			clientErr, ok := err.(*client.RequestError)
			if ok && clientErr.RateLimitExceeded() {
				return // Halt processing loop
			}
		}
	}
}

const FormatPng = "png"
const FormatZip = "zip"
const MimeZip = "application/zip"

func (is *ImageSettings) upgradePngToZipFormat() {
	// Save network bandwidth by requesting ZIP format (output will still be a PNG)
	if is.Format == FormatPng {
		is.Format = FormatZip
	}
}

func (p Processor) processFile(inputPath string, outputPath string, imageSettings ImageSettings) error {
	params := imageSettingsToParams(imageSettings)
	processedBytes, contentType, err := p.Client.RemoveFromFile(inputPath, p.APIKey, params)
	if err != nil {
		return err
	}

	if strings.Contains(contentType, MimeZip) {
		return p.processCompositeFile(outputPath, processedBytes)
	} else {
		return p.Storage.Write(outputPath, processedBytes)
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

	if len(imageSettings.BgImageFile) > 0 {
		params["bg_image_file"] = imageSettings.BgImageFile
	}

	if len(imageSettings.Format) > 0 {
		params["format"] = imageSettings.Format
	}

	if len(imageSettings.ExtraApiOptions) > 0 {
		values, err := url.ParseQuery(imageSettings.ExtraApiOptions)

		if err == nil {
			for key := range values {
				params[key] = values.Get(key)
			}
		} else {
			fmt.Printf("Unable to parse extra api options: %s\n", err)
		}
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

func (p Processor) processCompositeFile(outputPath string, processedBytes []byte) error {
	file, err := ioutil.TempFile("", "removebg.*.zip")
	if err != nil {
		return err
	}

	defer os.Remove(file.Name())

	_, err = file.Write(processedBytes)
	if err != nil {
		return err
	}

	// Convert output/foo.zip -> output/foo.png
	pngOutputPath := strings.TrimSuffix(outputPath, filepath.Ext(outputPath)) + ".png"

	return p.Compositor.Process(file.Name(), pngOutputPath)
}
