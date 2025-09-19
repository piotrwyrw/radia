package context

import (
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/piotrwyrw/otherproj/internal/types"
	"github.com/piotrwyrw/radia/radia/rtypes"
)

type Settings struct {
	ImageWidth  int
	ImageHeight int
}

type Context struct {
	RenderProgress binding.Float

	StatusLabel *widget.Label
	StatusText  binding.String

	PreviewImage *types.PreviewImage
	IsRendering  bool

	CurrentScene rtypes.Scene

	Settings Settings
}
