package cli

import (
	"github.com/remove-bg/go/client"
	"github.com/remove-bg/go/processor"
	"github.com/urfave/cli"
	"net/http"
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

		if len(outputDirectory) == 0 {
			return cli.NewExitError("Please specify the output directory", 1)
		}

		p := processor.Processor{
			APIKey: apiKey,
			Client: client.Client{
				HTTPClient: http.Client{},
			},
			FileWriter: processor.FileWriter{},
		}

		p.Process(inputPaths, outputDirectory)

		return nil
	}

	return app
}
