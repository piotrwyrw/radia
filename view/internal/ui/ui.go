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

	toolbar := widget.NewToolbar(widget.NewToolbarAction(theme.MediaPlayIcon(), func() {
		if ctx.IsRendering {
			return
		}

		go radia.InvokeRenderer(ctx, int32(imageWidth), int32(imageHeight))
	}))

	return container.NewBorder(
		container.NewVBox(
			toolbar,
			createSeparatorLine(),
		),
		nil,
		nil,
		nil,
		container.NewHSplit(
			container.NewVBox(),
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

	ctx := &context.Context{RenderProgress: binding.NewFloat(), StatusText: binding.NewString()}

	w := a.NewWindow("Radia")
	w.Resize(fyne.NewSize(1500, 900))

	split := createMainUI(ctx, imageWidth, imageHeight)

	ctx.StatusText.Set("Ready")

	w.SetContent(split)
	w.ShowAndRun()
}
