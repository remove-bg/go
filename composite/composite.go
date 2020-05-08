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

type Composite struct {
	Storage storage.StorageInterface
}

type imageDecoder = func(io.Reader) (image.Image, error)

func New() Composite {
	return Composite{
		Storage: storage.FileStorage{},
	}
}

func (c Composite) Process(inputZipPath string, outputImagePath string) error {
	fmt.Println("Extracting...")
	rgb, alpha, _ := extractImagesFromZip(inputZipPath)

	fmt.Println("Compositing...")
	composited := composite(rgb, alpha)

	fmt.Println("Saving...")
	c.savePng(composited, outputImagePath)

	return nil
}

func (c Composite) savePng(image *image.NRGBA, outputPath string) {
	buf := new(bytes.Buffer)
	png.Encode(buf, image)
	c.Storage.Write(outputPath, buf.Bytes())
}

const zipColorImageFileName = "color.jpg"
const zipAlphaImageFileName = "alpha.png"

func extractImagesFromZip(filename string) (rgb image.Image, alpha image.Image, err error) {
	archive, err := zip.OpenReader(filename)
	defer archive.Close()

	rgb, err = decodeZipImage(archive, zipColorImageFileName, jpeg.Decode)

	if err != nil {
		return nil, nil, err
	}

	alpha, err = decodeZipImage(archive, zipAlphaImageFileName, png.Decode)

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
