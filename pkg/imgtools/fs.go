package imgtools

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

func SavePNG(img image.Image, path string, name string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		if err := os.Mkdir(path, 0777); err != nil {
			return err
		}
	}

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
