package types

import (
	"image"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"github.com/piotrwyrw/otherproj/internal/util"
	"github.com/piotrwyrw/radia/radia/rimg"
)

type PreviewImage struct {
	ImageBuffer *image.RGBA
	ImageWidget *canvas.Image
}

func NewPreviewImage(buffer *image.RGBA, img *canvas.Image) *PreviewImage {
	return &PreviewImage{
		ImageBuffer: buffer,
		ImageWidget: img,
	}
}

func (img *PreviewImage) Update(data *rimg.Raster) {
	util.UpdateFyneImageWithRaster(data, img.ImageBuffer)
	fyne.Do(func() {
		img.ImageWidget.Refresh()
	})
}
