package processor

import (
	"path"
	"path/filepath"
	"strings"
)

const outputExtension = ".png"

func DetermineOutputPath(inputPath string, settings Settings) string {
	outputDirectory := settings.OutputDirectory
	inputDirectory, fileName := filepath.Split(inputPath)
	extentionlessFileName := strings.TrimSuffix(fileName, path.Ext(fileName))

	if len(outputDirectory) == 0 {
		return filepath.Join(inputDirectory, extentionlessFileName+"-removebg"+outputExtension)
	}

	return filepath.Join(outputDirectory, extentionlessFileName+outputExtension)
}
