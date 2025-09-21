package ui

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	ftheme "fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/piotrwyrw/otherproj/internal/cfgui"
	"github.com/piotrwyrw/otherproj/internal/radia"
	"github.com/piotrwyrw/otherproj/internal/reactive"
	"github.com/piotrwyrw/otherproj/internal/state"
	"github.com/piotrwyrw/otherproj/internal/ui/vtheme"
	"github.com/piotrwyrw/radia/radia/rscene"
)

func createSeparatorLine() *canvas.Line {
	sepLine := canvas.NewLine(color.Gray{Y: 15})
	sepLine.StrokeWidth = 1
	return sepLine
}

func createStatusBar(state *state.State) *fyne.Container {
	status := widget.NewLabelWithData(state.StatusText)
	status.TextStyle.Monospace = true
	status.Alignment = fyne.TextAlignLeading
	state.StatusLabel = status

	// Rendering progress indicator
	rect := canvas.NewRectangle(color.Transparent)
	rect.SetMinSize(fyne.NewSize(400, 0))
	progress := widget.NewProgressBarWithData(state.RenderProgress)
	progressContainer := container.NewStack(rect, progress)

	progress.Hide()
	state.IsRendering.ObserveFyne(func(newValue bool) {
		if newValue {
			progress.Show()
		} else {
			progress.Hide()
		}
	})

	statusBar := container.NewVBox(
		createSeparatorLine(),
		container.NewHBox(status, layout.NewSpacer(), progressContainer),
	)

	return statusBar
}

func createObjectList(state *state.State) fyne.CanvasObject {
	objectList := widget.NewList(
		func() int {
			return len(state.Context.CurrentScene.Objects)
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
			objects[2].(*widget.Label).SetText(state.Context.CurrentScene.Objects[i].Type)
			objects[5].(*widget.Label).SetText(state.Context.CurrentScene.Objects[i].Type)
		})

	state.SceneChanged.Observe(func() {
		objectList.Refresh()
	})

	return objectList
}

func createSidebar(state *state.State) (fyne.CanvasObject, error) {
	settings, panel, err := cfgui.CreateSettingsPanel(&state.Context)
	if err != nil {
		return nil, err
	}

	state.Settings = panel

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("Settings", ftheme.SettingsIcon(), container.NewVScroll(settings)),
		container.NewTabItemWithIcon("Objects", ftheme.StorageIcon(), createObjectList(state)),
	)

	return tabs, nil
}

func createImagePreview(state *state.State) fyne.CanvasObject {
	state.PreviewImage.Create(state.Context.Settings.ImageWidth, state.Context.Settings.ImageHeight)
	return state.PreviewImage.ImageWidget
}

func createToolbar(state *state.State) fyne.CanvasObject {
	renderAction := widget.NewToolbarAction(ftheme.MediaPlayIcon(), func() {
		if state.IsRendering.Get() {
			return
		}

		go radia.InvokeRenderer(state)
	})

	openAction := widget.NewToolbarAction(ftheme.FolderOpenIcon(), func() {
		go showSceneOpenDialog(state)
	})

	saveAction := widget.NewToolbarAction(ftheme.DocumentSaveIcon(), func() {
		go showSceneSaveDialog(state)
	})

	saveImageAction := widget.NewToolbarAction(ftheme.FileImageIcon(), func() {
		go showImageSaveDialog(state)
	})

	actions := []*widget.ToolbarAction{renderAction, openAction, saveAction, saveImageAction}

	state.IsRendering.ObserveFyne(func(newValue bool) {
		for _, action := range actions {
			if newValue {
				action.Disable()
			} else {
				action.Enable()
			}
		}
	})

	toolbar := widget.NewToolbar(renderAction, openAction, saveAction, saveImageAction)

	return toolbar
}

func createMainUI(state *state.State) (fyne.CanvasObject, error) {
	props, err := createSidebar(state)
	if err != nil {
		return nil, err
	}

	toolbar := createToolbar(state)
	statusBar := createStatusBar(state)
	preview := createImagePreview(state)

	var borderContainer *fyne.Container

	borderContainer = container.NewBorder(
		nil,
		statusBar,
		nil,
		nil,
		preview,
	)

	// TODO Replace with ObservableValue
	state.PreviewImage.OnUpdate = func(img *canvas.Image) {
		borderContainer.Objects[0] = img
	}

	mainSplit := container.NewHSplit(
		props,
		borderContainer,
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

	a := app.NewWithID("view.master")

	theme := vtheme.RadiaTheme{
		Fallback: a.Settings().Theme(),
	}

	a.Settings().SetTheme(theme)

	w := a.NewWindow("Radia Studio")
	w.Resize(fyne.NewSize(1500, 900))

	s := &state.State{
		MainWindow:     w,
		RenderProgress: binding.NewFloat(),
		StatusText:     binding.NewString(),
		SceneChanged:   reactive.NewSignal(),
		Context: state.RenderContext{
			Settings: state.RenderSettings{
				Samples:     10,
				MaxBounces:  100,
				Threads:     0,
				ImageWidth:  imageWidth,
				ImageHeight: imageHeight,
			},
		},
		IsRendering: reactive.NewObservableValue(false),
	}
	s.Context.CurrentScene = *rscene.NewBlankScene()

	ui, err := createMainUI(s)
	if err != nil {
		return err
	}

	s.StatusText.Set("Ready")

	w.SetContent(ui)
	w.Show()
	a.Run()

	return nil
}
