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
	ftheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/piotrwyrw/otherproj/internal/context"
	"github.com/piotrwyrw/otherproj/internal/radia"
	"github.com/piotrwyrw/otherproj/internal/types"
	"github.com/piotrwyrw/otherproj/internal/ui/vtheme"
	"github.com/piotrwyrw/radia/radia/rscene"
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

func createObjectList(ctx *context.Context) fyne.CanvasObject {
	objectList := widget.NewList(
		func() int {
			return len(ctx.CurrentScene.Objects)
		},
		func() fyne.CanvasObject {
			typeLabel := widget.NewLabel("")
			typeLabel.Alignment = fyne.TextAlignLeading

			nameLabel := widget.NewLabel("")
			nameLabel.Alignment = fyne.TextAlignTrailing

			return container.NewHBox(
				layout.NewSpacer(),
				widget.NewIcon(ftheme.FileIcon()),
				nameLabel,
				layout.NewSpacer(),
				widget.NewIcon(ftheme.GridIcon()),
				typeLabel,
				layout.NewSpacer(),
			)
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			objects := o.(*fyne.Container).Objects
			objects[2].(*widget.Label).SetText(ctx.CurrentScene.Objects[i].Type)
			objects[5].(*widget.Label).SetText(ctx.CurrentScene.Objects[i].Type)
		})

	return objectList
}

func createSidebar(ctx *context.Context) (fyne.CanvasObject, error) {
	settings, err := createSettingsPanel(&ctx.CurrentScene)
	if err != nil {
		return nil, err
	}

	return container.NewAppTabs(
		container.NewTabItemWithIcon("Settings", ftheme.SettingsIcon(), container.NewVScroll(settings)),
		container.NewTabItemWithIcon("Outliner", ftheme.StorageIcon(), createObjectList(ctx)),
	), nil
}

func createImagePreview(ctx *context.Context) fyne.CanvasObject {
	buff := image.NewRGBA(image.Rect(0, 0, ctx.Settings.ImageWidth, ctx.Settings.ImageHeight))
	img := canvas.NewImageFromImage(buff)
	img.FillMode = canvas.ImageFillContain

	for y := 0; y < ctx.Settings.ImageHeight; y++ {
		for x := 0; x < ctx.Settings.ImageWidth; x++ {
			buff.Set(x, y, color.Gray{Y: 5})
		}
	}

	ctx.PreviewImage = types.NewPreviewImage(buff, img)

	return img
}

func createToolbar(ctx *context.Context) fyne.CanvasObject {
	toolbar := widget.NewToolbar(
		widget.NewToolbarAction(ftheme.MediaPlayIcon(), func() {
			if ctx.IsRendering {
				return
			}

			go radia.InvokeRenderer(ctx)
		}),
		widget.NewToolbarAction(ftheme.FolderOpenIcon(), func() {

		}),
	)

	return toolbar
}

func createMainUI(ctx *context.Context) (fyne.CanvasObject, error) {
	props, err := createSidebar(ctx)
	if err != nil {
		return nil, err
	}

	toolbar := createToolbar(ctx)
	statusBar := createStatusBar(ctx)
	preview := createImagePreview(ctx)

	mainSplit := container.NewHSplit(
		props,
		container.NewBorder(
			nil,
			statusBar,
			nil,
			nil,
			preview,
		),
	)
	mainSplit.SetOffset(0.25)

	return container.NewBorder(
		container.NewVBox(
			toolbar,
			createSeparatorLine(),
		),
		nil,
		nil,
		nil,
		mainSplit,
	), nil
}

func CreateUI() error {
	imageWidth, imageHeight := 1500, 900

	a := app.New()
	a.Settings().SetTheme(vtheme.RadiaTheme{})

	ctx := &context.Context{RenderProgress: binding.NewFloat(), StatusText: binding.NewString(), Settings: context.Settings{
		ImageWidth:  imageWidth,
		ImageHeight: imageHeight,
	}}
	ctx.CurrentScene = *rscene.NewBlankScene()

	w := a.NewWindow("Radia Studio")
	w.Resize(fyne.NewSize(1500, 900))

	ui, err := createMainUI(ctx)
	if err != nil {
		return err
	}

	ctx.StatusText.Set("Ready")

	w.SetContent(ui)
	w.Show()
	a.Run()

	return nil
}
