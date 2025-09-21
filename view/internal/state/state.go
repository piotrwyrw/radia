package state

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/piotrwyrw/otherproj/internal/cfgui"
	"github.com/piotrwyrw/otherproj/internal/reactive"
	"github.com/piotrwyrw/otherproj/internal/types"
	"github.com/piotrwyrw/radia/radia/rtypes"
)

type RenderSettings struct {
	Samples     int `ui:"Samples"`
	MaxBounces  int `ui:"Max Bounces"`
	Threads     int `ui:"Thread Limit"`
	ImageWidth  int `ui:"Image Width"`
	ImageHeight int `ui:"Image Height"`
}

type RenderContext struct {
	CurrentScene rtypes.Scene   `ui:"Scene"`
	Settings     RenderSettings `ui:"Renderer Settings"`
}

type State struct {
	MainWindow fyne.Window

	RenderProgress binding.Float

	StatusLabel *widget.Label
	StatusText  binding.String

	PreviewImage types.PreviewImage
	IsRendering  *reactive.ObservableValue[bool]

	SceneChanged *reactive.Signal
	Context      RenderContext
	Settings     *cfgui.SettingsPanel
}
