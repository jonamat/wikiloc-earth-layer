package imgtools

import (
	"bytes"
	"image"

	"github.com/srwiley/oksvg"
	"github.com/srwiley/rasterx"
)

func RawSVGToPNGImage(svg []byte, size int) (image.Image, error) {
	reader := bytes.NewReader(svg)

	icon, err := oksvg.ReadIconStream(reader)
	if err != nil {
		return nil, err
	}
	icon.SetTarget(0, 0, float64(size), float64(size))

	rgba := image.NewCMYK(image.Rect(0, 0, size, size))
	scanner := rasterx.NewScannerGV(size, size, rgba, rgba.Bounds())
	// TODO add precision to svg conversion

	icon.Draw(rasterx.NewDasher(size, size, scanner), 1)

	img := rgba.SubImage(image.Rect(0, 0, size, size))

	return img, nil
}
