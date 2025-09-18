package context

import (
	"image"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type Context struct {
	RenderProgress binding.Float

	StatusLabel *widget.Label
	StatusText  binding.String

	RenderOutputBuffer *image.RGBA
	RenderOutputImage  *canvas.Image
	IsRendering        bool
}
