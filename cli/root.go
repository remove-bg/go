package cli

import (
	"fmt"

	"github.com/urfave/cli"
)

// Bootstrap the CLI
func Bootstrap() *cli.App {
	app := cli.NewApp()

	app.Action = func(c *cli.Context) error {
		fmt.Println("Hello World")
		return nil
	}

	return app
}
