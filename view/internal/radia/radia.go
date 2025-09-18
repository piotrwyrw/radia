package radia

import (
	"math"

	"fyne.io/fyne/v2"
	"github.com/piotrwyrw/otherproj/internal/context"
	"github.com/piotrwyrw/otherproj/internal/util"
	"github.com/piotrwyrw/radia/radia/rcolor"
	"github.com/piotrwyrw/radia/radia/rimg"
	"github.com/piotrwyrw/radia/radia/rmaterial"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/robject"
	"github.com/piotrwyrw/radia/radia/rscene"
	"github.com/piotrwyrw/radia/radia/rshapes"
	"github.com/piotrwyrw/radia/radia/rtracer"
	"github.com/piotrwyrw/radia/radia/rtypes"
	"github.com/sirupsen/logrus"
)

func InvokeRenderer(ctx *context.Context, imageWidth int32, imageHeight int32) {
	ctx.StatusText.Set("Rendering ...")
	ctx.IsRendering = true
	ctx.RenderProgress.Set(0.0)

	env, err := rimg.RasterFromPNG("world.png")
	if err != nil {
		logrus.Errorf("Could not render: %v\n", err)
		ctx.IsRendering = false
		return
	}

	mat := rmaterial.UniversalMaterial{
		Color:     rcolor.Color{R: 0.1, G: 0.6, B: 0.3},
		Emission:  rcolor.ColorBlack(),
		Roughness: 0.5,
	}

	scene := &rtypes.Scene{
		Objects: []rtypes.ShapeWrapper{
			robject.WrapShape(&rshapes.Sphere{
				Center:   rmath.Vec3d{X: 0.0, Y: 0.0, Z: 2.0},
				Radius:   0.2,
				Material: robject.WrapShapeMaterial(&mat),
			}),
			robject.WrapShape(&rshapes.Sphere{
				Center: rmath.Vec3d{X: -0.4, Y: 0.0, Z: 2.0},
				Radius: 0.1,
				Material: robject.WrapShapeMaterial(&rmaterial.UniversalMaterial{
					Color:      rcolor.ColorWhite(),
					Emission:   rcolor.ColorWhite(),
					Brightness: 3.0,
					Roughness:  10.0,
				}),
			}),
		},
		Camera: rtypes.Camera{
			Location:    rmath.Vec3d{X: 0.0, Y: 0.0, Z: 0.0},
			Facing:      rmath.Vec3d{X: 0.0, Y: 0.0, Z: 1.0},
			FocalLength: 1.0,
		},
		WorldMat: robject.WrapEnvironmentMaterial(&rmaterial.Sky{
			Image:         env,
			IOR:           1.0,
			FallbackColor: rcolor.ColorBlack(),
			Intensity:     1.2,
			Azimuth:       math.Pi * 1.5,
			Elevation:     0.0,
		}),
	}

	_ = rscene.SaveSceneJSON(scene, "scene_js.json")
	_, err = rscene.LoadSceneJSON("scene_js.json")
	if err != nil {
		logrus.Errorf("Could not load scene json: %v\n", err)
	}

	destination := rimg.NewRaster(imageWidth, imageHeight)

	var lastPercent float64 = 0.0
	rtracer.TraceAllRays(scene, destination, 50, 100, func(x, y, n int32) {
		progress := float64(n) / float64(imageWidth*imageHeight)
		if int(progress*100) != int(lastPercent*100) {
			lastPercent = progress
			ctx.RenderProgress.Set(progress)
		}
	})

	util.UpdateFyneImageWithRaster(destination, ctx.RenderOutputBuffer)

	fyne.DoAndWait(func() {
		ctx.StatusText.Set("Ready")
		ctx.RenderOutputImage.Refresh()
		ctx.IsRendering = false
		ctx.RenderProgress.Set(0.0)
	})
}
