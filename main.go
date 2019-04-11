package main

import (
	"log"
	"os"

	"github.com/remove-bg/go/cli"
)

func main() {
	app := cli.Bootstrap()
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
