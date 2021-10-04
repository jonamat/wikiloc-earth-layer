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

	dasher := rasterx.NewDasher(size, size, rasterx.NewScannerGV(size, size, rgba, rgba.Bounds()))
	icon.Draw(dasher, 1)

	img := rgba.SubImage(image.Rect(0, 0, size, size))

	return img, nil
}
