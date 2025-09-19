package types

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/piotrwyrw/otherproj/internal/util"
	"github.com/piotrwyrw/radia/radia/rimg"
)

type PreviewImage struct {
	ImageBuffer *image.RGBA
	ImageWidget *canvas.Image
	OnUpdate    func(img *canvas.Image)
}

func (img *PreviewImage) Create(width, height int) {
	img.ImageBuffer = image.NewRGBA(image.Rect(0, 0, width, height))

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.ImageBuffer.Set(x, y, color.Gray{Y: 30})
		}
	}

	img.ImageWidget = canvas.NewImageFromImage(img.ImageBuffer)
	img.ImageWidget.FillMode = canvas.ImageFillContain

	if img.OnUpdate != nil {
		img.OnUpdate(img.ImageWidget)
	}
}

func (img *PreviewImage) Refresh() {
	fyne.Do(func() {
		img.ImageWidget.Refresh()
	})
}

func (img *PreviewImage) Update(data *rimg.Raster) {
	util.UpdateFyneImageWithRaster(data, img.ImageBuffer)
	img.Refresh()
}
