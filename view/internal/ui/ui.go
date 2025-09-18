package ui

import (
	"image"
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	theme2 "fyne.io/x/fyne/theme"
	"github.com/piotrwyrw/otherproj/internal/context"
	"github.com/piotrwyrw/otherproj/internal/radia"
)

func createSeparatorLine() *canvas.Line {
	sepLine := canvas.NewLine(color.Gray{Y: 15})
	sepLine.StrokeWidth = 1
	return sepLine
}

func createStatusBar(ctx *context.Context) *fyne.Container {
	status := widget.NewLabelWithData(ctx.StatusText)
	status.TextStyle.Monospace = true
	status.Alignment = fyne.TextAlignLeading
	ctx.StatusLabel = status

	// Rendering progress indicator
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(400, 0))
	progress := widget.NewProgressBarWithData(ctx.RenderProgress)
	progressContainer := container.NewStack(rect, progress)

	statusBar := container.NewVBox(
		createSeparatorLine(),
		container.NewHBox(status, layout.NewSpacer(), progressContainer),
	)

	return statusBar
}

func createPropertiesTab() fyne.CanvasObject {
	return container.NewAppTabs(
		container.NewTabItemWithIcon("Settings", theme.SettingsIcon(), widget.NewLabel("Settings")),
		container.NewTabItemWithIcon("Object", theme.InfoIcon(), widget.NewLabel("Object Settings")),
	)
}

func createMainUI(ctx *context.Context, imageWidth, imageHeight int) fyne.CanvasObject {
	buff := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	img := canvas.NewImageFromImage(buff)
	img.FillMode = canvas.ImageFillContain

	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			buff.Set(x, y, color.Gray{Y: 30})
		}
	}

	ctx.RenderOutputImage = img
	ctx.RenderOutputBuffer = buff

	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
			if ctx.IsRendering {
				return
			}

			go radia.InvokeRenderer(ctx, int32(imageWidth), int32(imageHeight))
		}),
		widget.NewToolbarAction(theme.FolderOpenIcon(), func() {

		}),
	)

	return container.NewBorder(
		container.NewVBox(
			toolbar,
			createSeparatorLine(),
		),
		nil,
		nil,
		nil,
		container.NewHSplit(
			createPropertiesTab(),
			container.NewBorder(
				nil,
				createStatusBar(ctx),
				nil,
				nil,
				img,
			),
		),
	)
}

func CreateUI() {
	imageWidth, imageHeight := 1500, 900

	a := app.New()
	a.Settings().SetTheme(theme2.AdwaitaTheme())

	ctx := &context.Context{RenderProgress: binding.NewFloat(), StatusText: binding.NewString()}

	w := a.NewWindow("Radia Studio")
	w.Resize(fyne.NewSize(1500, 900))

	ui := createMainUI(ctx, imageWidth, imageHeight)

	ctx.StatusText.Set("Ready")

	w.SetContent(ui)
	w.Show()
	a.Run()
}
