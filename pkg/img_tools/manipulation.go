package imgtools

import (
	"image"
	"image/color"
	"image/draw"
)

type circle struct {
	p image.Point
	r int
}

func (c *circle) ColorModel() color.Model {
	return color.AlphaModel
}

func (c *circle) Bounds() image.Rectangle {
	return image.Rect(c.p.X-c.r, c.p.Y-c.r, c.p.X+c.r, c.p.Y+c.r)
}

func (c *circle) At(x, y int) color.Color {
	_x, _y, _r := float64(x-c.p.X)+0.5, float64(y-c.p.Y)+0.5, float64(c.r)
	if _x*_x+_y*_y < _r*_r {
		return color.Alpha{255}
	}
	return color.Alpha{0}
}

func CropIntoCirlce(src *image.Image, radius int) *image.RGBA {
	mask := image.NewRGBA(image.Rect(0, 0, radius*2, radius*2))

	// Apply transparent bg
	draw.Draw(mask, mask.Bounds(), image.Transparent, image.Point{}, draw.Src)

	// Create a circle
	circle := &circle{image.Point{radius, radius}, radius}

	// Apply mask to source
	draw.DrawMask(mask, mask.Bounds(), *src, image.Point{}, circle, image.Point{}, draw.Over)

	return mask
}
