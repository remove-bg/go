package cli

import (
	"fmt"

	"github.com/urfave/cli"
)

// Bootstrap the CLI
func Bootstrap() *cli.App {
	app := cli.NewApp()

	app.Name = "removebg"
	app.Description = "Remove image background - 100% automatically"
	app.Version = "0.1.0"

	app.Action = func(c *cli.Context) error {
		fmt.Println("Hello World")
		return nil
	}

	return app
}
