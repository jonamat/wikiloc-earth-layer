package imgtools

import (
	"image"

	vp "github.com/spf13/viper"
)

var size = vp.GetInt("iconSidePx")

func MakeIcon(RawSVG *[]byte) (*image.RGBA, error) {
	// Parse rawSVG and create a grayscale PNG
	fullIcon, err := RawSVGToPNGImage(*RawSVG, size)
	if err != nil {
		return nil, err
	}

	// Cut out the full icon into a circle with transparent background
	croppedIcon := CropIntoCirlce(&fullIcon, size/2)

	return croppedIcon, nil
}
