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
	"io"
)

//go:generate counterfeiter . CompositorInterface
type CompositorInterface interface {
	Process(inputZipPath string, outputImagePath string) error
}

type Compositor struct {
	Storage storage.StorageInterface
}

type imageDecoder = func(io.Reader) (image.Image, error)

func New() Compositor {
	return Compositor{
		Storage: storage.FileStorage{},
	}
}

func (c Compositor) Process(inputZipPath string, outputImagePath string) error {
	if !c.Storage.FileExists(inputZipPath) {
		return fmt.Errorf("Could not locate zip: %s", inputZipPath)
	}

	rgb, alpha, err := extractImagesFromZip(inputZipPath)

	if err != nil {
		return err
	}

	composited := composite(rgb, alpha)

	c.savePng(composited, outputImagePath)

	return nil
}

func (c Compositor) savePng(image *image.NRGBA, outputPath string) {
	buf := new(bytes.Buffer)
	png.Encode(buf, image)
	c.Storage.Write(outputPath, buf.Bytes())
}

const zipColorImageFileName = "color.jpg"
const zipAlphaImageFileName = "alpha.png"

func extractImagesFromZip(filename string) (rgb image.Image, alpha image.Image, err error) {
	archive, err := zip.OpenReader(filename)
	if err != nil {
		return nil, nil, err
	}

	defer archive.Close()

	alpha, err = decodeZipImage(archive, zipAlphaImageFileName, png.Decode)
	if err != nil {
		return nil, nil, err
	}

	rgb, err = decodeZipImage(archive, zipColorImageFileName, jpeg.Decode)
	if err != nil {
		return nil, nil, err
	}

	return rgb, alpha, nil
}

func decodeZipImage(archive *zip.ReadCloser, fileName string, decoder imageDecoder) (image.Image, error) {
	for _, f := range archive.File {
		if f.Name == fileName {
			rc, err := f.Open()
			defer rc.Close()

			if err != nil {
				return nil, err
			}

			return decoder(rc)
		}
	}

	return nil, fmt.Errorf("Unable to find image in ZIP: %s", fileName)
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
