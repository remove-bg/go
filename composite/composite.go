package composite

import (
	"archive/zip"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
)

func Process(inputPath string, outputPath string) error {
	fmt.Println("Extracting...")
	rgb, alpha := readZIP(inputPath)

	fmt.Println("Compositing...")
	composited := composite(rgb, alpha)

	fmt.Println("Saving...")
	savePNG(composited, outputPath)

	return nil
}

func readZIP(filename string) (image.Image, image.Image) {
	r, err := zip.OpenReader(filename)
	defer r.Close()

	if err != nil {
		// TODO
	}

	var rgb image.Image
	var alpha image.Image

	for _, f := range r.File {
		if f.Name == "color.jpg" {
			rc, err := f.Open()
			if err != nil {
				// TODO
			}

			rgb, err = jpeg.Decode(rc)

			if err != nil {
				// TODO
			}

			rc.Close()
		}

		if f.Name == "alpha.png" {
			rc, err := f.Open()
			if err != nil {
				// TODO
			}

			alpha, err = png.Decode(rc)
			if err != nil {
				// TODO
			}

			rc.Close()
		}
	}

	return rgb, alpha
}

func composite(rgb image.Image, alpha image.Image) *image.NRGBA {
	dimensions := rgb.Bounds().Max
	width := dimensions.X
	height := dimensions.Y

	composited := image.NewNRGBA(image.Rect(0, 0, width, height))

	colorModel := composited.ColorModel()

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			rgbColor := (colorModel.Convert(rgb.At(x, y))).(color.NRGBA)
			alphaColor := (alpha.At(x, y)).(color.Gray)
			rgbColor.A = alphaColor.Y

			composited.SetNRGBA(x, y, rgbColor)
		}
	}

	return composited
}

func savePNG(image *image.NRGBA, filename string) {
	outputFile, err := os.Create(filename)
	defer outputFile.Close()

	if err != nil {
		// TODO
	}
	png.Encode(outputFile, image)
}
