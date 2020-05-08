package composite

import (
	"github.com/remove-bg/go/storage"

	"archive/zip"
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
)

type Composite struct {
	Storage storage.StorageInterface
}

func New() Composite {
	return Composite{
		Storage: storage.FileStorage{},
	}
}

func (c Composite) Process(inputPath string, outputPath string) error {
	fmt.Println("Extracting...")
	rgb, alpha := readZip(inputPath)

	fmt.Println("Compositing...")
	composited := composite(rgb, alpha)

	fmt.Println("Saving...")
	c.savePng(composited, outputPath)

	return nil
}

func (c Composite) savePng(image *image.NRGBA, outputPath string) {
	buf := new(bytes.Buffer)
	png.Encode(buf, image)
	c.Storage.Write(outputPath, buf.Bytes())
}

func readZip(filename string) (rgb image.Image, alpha image.Image) {
	r, err := zip.OpenReader(filename)
	defer r.Close()

	if err != nil {
		// TODO
	}

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
