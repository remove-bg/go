package cli

import (
	"github.com/remove-bg/go/processor"
	"github.com/urfave/cli"
)

// Bootstrap the CLI
func Bootstrap() *cli.App {
	app := cli.NewApp()

	app.Name = "removebg"
	app.Description = "Remove image background - 100% automatically"
	app.Version = "0.1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "api-key",
			Usage:  "API key",
			EnvVar: "REMOVE_BG_API_KEY",
		},
		cli.StringFlag{
			Name:  "output-directory",
			Usage: "Output directory",
		},
		cli.StringFlag{
			Name:  "size",
			Usage: "Image size",
		},
		cli.StringFlag{
			Name:  "type",
			Usage: "Image type",
		},
		cli.StringFlag{
			Name:  "channels",
			Usage: "Image channels",
		},
		cli.StringFlag{
			Name:  "bg-color",
			Usage: "Image background color",
		},
	}

	app.Action = func(c *cli.Context) error {
		apiKey := c.String("api-key")
		outputDirectory := c.String("output-directory")
		inputPaths := c.Args()

		if len(apiKey) == 0 {
			return cli.NewExitError("API key must be specified", 1)
		}

		if len(inputPaths) == 0 {
			return cli.NewExitError("Please specify one or more files", 1)
		}

		p := processor.NewProcessor(apiKey)
		s := processor.Settings{
			OutputDirectory: outputDirectory,
			ImageSettings: processor.ImageSettings{
				Size:     c.String("size"),
				Type:     c.String("type"),
				Channels: c.String("channels"),
				BgColor:  c.String("bg-color"),
			},
		}

		p.Process(inputPaths, s)

		return nil
	}

	return app
}
