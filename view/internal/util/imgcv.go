package util

import (
	"image"
	"image/color"

	"github.com/piotrwyrw/radia/radia/rimg"
)

func UpdateFyneImageWithRaster(raster *rimg.Raster, destination *image.RGBA) {
	for y := 0; y < int(raster.Height); y++ {
		for x := 0; x < int(raster.Width); x++ {
			rt := raster.Get(int32(x), int32(y))
			r, g, b := rt.SDLColor()

			destination.Set(x, y, color.RGBA{R: r, G: g, B: b, A: 255})
		}
	}
}
