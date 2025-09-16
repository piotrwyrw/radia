package context

import (
	"image"

	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/data/binding"
)

type Context struct {
	RenderProgress     binding.Float
	StatusText         binding.String
	RenderOutputBuffer *image.RGBA
	RenderOutputImage  *canvas.Image
	IsRendering        bool
}
