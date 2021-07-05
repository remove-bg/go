package cmd

import (
	"./composite"
	"github.com/spf13/cobra"
	"log"
)

var zip2pngCmd = &cobra.Command{
	Short: "Converts a remove.bg ZIP to a PNG",
	Use:   "zip2png <input.zip> <output_path.png>",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		inputZipPath := args[0]
		outputImagePath := args[1]
		composite := composite.New()

		err := composite.Process(inputZipPath, outputImagePath)

		if err != nil {
			return err
		}

		log.Printf("Processed zip: %s -> %s\n", inputZipPath, outputImagePath)
		return nil
	},
}

func init() {
	RootCmd.AddCommand(zip2pngCmd)
}
