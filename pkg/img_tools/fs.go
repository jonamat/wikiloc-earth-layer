package imgtools

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func SavePNG(img image.Image, name string, path string) error {
	file, err := os.Create(fmt.Sprintf("%s%s.png", path, name))
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}
