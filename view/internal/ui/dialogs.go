package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/piotrwyrw/otherproj/internal/errh"
	"github.com/piotrwyrw/otherproj/internal/state"
	"github.com/piotrwyrw/radia/radia/rparser"
	"github.com/piotrwyrw/radia/radia/rscene"
	"github.com/sirupsen/logrus"
)

func showSceneOpenDialog(state *state.State) {
	window := state.MainWindow

	d := dialog.NewFileOpen(func(read fyne.URIReadCloser, e error) {
		if e != nil {
			logrus.Errorf("Could not open scene file: %v\n", e)
			return
		}

		if read == nil {
			return
		}

		scene, err := rparser.LoadSceneJSON(read.URI().Path())
		if err != nil {
			errh.DisplayError(state, "Could Not Load Scene", err)
			logrus.Errorf("Could not load scene file: %v\n", err)
			return
		}
		state.SceneChanged.Notify()
		state.Context.CurrentScene = *scene
		state.Settings.SetDefaultValues()
	}, window)
	d.Resize(window.Canvas().Size())
	fyne.Do(func() {
		d.Show()
	})
}

func showSceneSaveDialog(state *state.State) {
	window := state.MainWindow
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
			errh.DisplayError(state, "Could Not Save Scene", err)
			logrus.Errorf("Could not save scene file: %v\n", err)
			return
		}
	}, window)
	d.Resize(window.Canvas().Size())
	fyne.Do(func() {
		d.Show()
	})
}

func showImageSaveDialog(state *state.State) {
	if state.PreviewImage.Raster == nil {
		return
	}
	
	window := state.MainWindow
	d := dialog.NewFileSave(func(write fyne.URIWriteCloser, e error) {
		if e != nil {
			logrus.Errorf("Could not save image: %v\n", e)
			return
		}

		if write == nil {
			return
		}

		err := state.PreviewImage.Raster.SavePPM(write.URI().Path())
		if err != nil {
			errh.DisplayError(state, "Could Not Save Image", err)
			logrus.Errorf("Could not save image: %v\n", err)
		}
	}, window)
	d.Resize(window.Canvas().Size())
	fyne.Do(func() {
		d.Show()
	})
}
