package cmd

import (
	"errors"
	"fmt"
	"github.com/remove-bg/go/processor"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

const defaultLargeBatchSize = 50

var (
	apiKey                    string
	confirmBatchOver          int
	outputDirectory           string
	reprocessExisting         bool
	skipPngFormatOptimization bool
	imageSize                 string
	imageType                 string
	imageFormat               string
	imageChannels             string
	bgColor                   string
	bgImageFile               string
	extraApiOptions           string
)

// RootCmd is the entry point of command-line execution
var RootCmd = &cobra.Command{
	Short: "Remove image background - 100% automatically",
	Use:   "removebg <file>...",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(apiKey) == 0 {
			return errors.New("API key must be specified")
		}

		if len(args) == 0 {
			return errors.New("please specify one or more files")
		}

		p := processor.NewProcessor(apiKey, cmd.Version)
		s := processor.Settings{
			OutputDirectory:            outputDirectory,
			ReprocessExisting:          reprocessExisting,
			SkipPngFormatOptimization:  skipPngFormatOptimization,
			LargeBatchConfirmThreshold: confirmBatchOver,
			ImageSettings: processor.ImageSettings{
				Size:            imageSize,
				Type:            imageType,
				Channels:        imageChannels,
				BgColor:         bgColor,
				BgImageFile:     bgImageFile,
				OutputFormat:    strings.ToLower(imageFormat),
				ExtraApiOptions: extraApiOptions,
			},
		}

		p.Process(args, s)

		return nil
	},
}

func ConfigureVersion(version string, commit string) {
	RootCmd.Version = version
	RootCmd.SetVersionTemplate(fmt.Sprintf("%s\n%s\n", version, commit))
}

func init() {
	RootCmd.Flags().StringVar(&apiKey, "api-key", "", "API key (required) or set REMOVE_BG_API_KEY environment variable")
	RootCmd.Flags().StringVar(&outputDirectory, "output-directory", "", "Output directory")
	RootCmd.Flags().BoolVar(&reprocessExisting, "reprocess-existing", false, "Reprocess and overwrite any already processed images")
	RootCmd.Flags().BoolVar(&skipPngFormatOptimization, "skip-png-format-optimization", false, "Skip optimizing PNG format as ZIP to save bandwidth (default false)")
	RootCmd.Flags().IntVar(&confirmBatchOver, "confirm-batch-over", defaultLargeBatchSize, "Confirm any batches over this size (-1 to disable)")
	RootCmd.Flags().StringVar(&imageSize, "size", "auto", "Image size")
	RootCmd.Flags().StringVar(&imageType, "type", "", "Image type")
	RootCmd.Flags().StringVar(&imageFormat, "format", "png", "Image format")
	RootCmd.Flags().StringVar(&imageChannels, "channels", "", "Image channels")
	RootCmd.Flags().StringVar(&bgColor, "bg-color", "", "Image background color")
	RootCmd.Flags().StringVar(&bgImageFile, "bg-image-file", "", "Adds a background image from a file")
	RootCmd.Flags().StringVar(&extraApiOptions, "extra-api-options", "", "Extra options to forward to the API (format: 'option1=val1&option2=val2')")

	if len(apiKey) == 0 {
		apiKey = os.Getenv("REMOVE_BG_API_KEY")
	}
}
