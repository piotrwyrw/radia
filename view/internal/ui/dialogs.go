package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/piotrwyrw/otherproj/internal/state"
	"github.com/piotrwyrw/radia/radia/rscene"
	"github.com/sirupsen/logrus"
)

func getCurrentWindow() fyne.Window {
	return fyne.CurrentApp().Driver().AllWindows()[0]
}

func showSceneOpenDialog(state *state.State) {
	window := getCurrentWindow()
	d := dialog.NewFileOpen(func(read fyne.URIReadCloser, e error) {
		if e != nil {
			logrus.Errorf("Could not open scene file: %v\n", e)
			return
		}

		if read == nil {
			return
		}

		scene, err := rscene.LoadSceneJSON(read.URI().Path())
		if err != nil {
			logrus.Errorf("Could not load scene file: %v\n", err)
			return
		}
		state.Context.CurrentScene = *scene
		state.Settings.SetDefaultValues()
	}, window)
	d.Resize(window.Canvas().Size())
	fyne.Do(func() {
		d.Show()
	})
}

func showSceneSaveDialog(state *state.State) {
	window := getCurrentWindow()
	d := dialog.NewFileSave(func(write fyne.URIWriteCloser, e error) {
		if e != nil {
			logrus.Errorf("Could not save scene file: %v\n", e)
			return
		}

		if write == nil {
			return
		}

		err := rscene.SaveSceneJSON(&state.Context.CurrentScene, write.URI().Path())
		if err != nil {
			logrus.Errorf("Could not save scene file: %v\n", err)
			return
		}
	}, window)
	d.Resize(window.Canvas().Size())
	fyne.Do(func() {
		d.Show()
	})
}
