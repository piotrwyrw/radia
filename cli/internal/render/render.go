package render

import (
	"cli/internal/parameters"

	"github.com/piotrwyrw/radia/radia/radia"
	"github.com/piotrwyrw/radia/radia/rimg"
	"github.com/piotrwyrw/radia/radia/rparser"
	"github.com/piotrwyrw/radia/radia/rtracer"
	"github.com/sirupsen/logrus"
)

func Scene(scenePath string, params *parameters.RenderParameters) error {
	radia.Initialize()
	logrus.Infof("Rendering scene \"%s\"", scenePath)
	scene, err := rparser.LoadSceneJSON(scenePath)
	if err != nil {
		logrus.Errorf("Error loading scene: %v", err)
		return err
	}
	img := rimg.NewRaster(int32(params.Width), int32(params.Height))
	var lastProgress = 0
	rtracer.TraceAllRays(scene, img, params.Samples, params.Bounces, params.Threads, func(n int32, progress float64) {
		if int(progress*100) <= lastProgress {
			return
		}

		lastProgress = int(progress * 100)
		logrus.Infof("Progress: %d%%", lastProgress)
	})
	logrus.Infof("Saving output")
	err = img.SavePPM(params.Output)
	if err != nil {
		logrus.Errorf("Error saving output: %v", err)
		return err
	}
	logrus.Infof("Done.")
	return nil
}
