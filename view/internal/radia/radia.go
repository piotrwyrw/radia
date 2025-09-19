package radia

import (
	"fmt"
	"math"
	"time"

	"github.com/piotrwyrw/otherproj/internal/context"
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

func InvokeRenderer(ctx *context.Context) {
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

	//scene := rscene.NewBlankScene()

	_ = rscene.SaveSceneJSON(scene, "scene_js.json")
	s, err := rscene.LoadSceneJSON("scene_js.json")
	if err != nil {
		logrus.Errorf("Could not load scene json: %v\n", err)
	}
	ctx.CurrentScene = *s

	imgWidth := int32(ctx.Settings.ImageWidth)
	imgHeight := int32(ctx.Settings.ImageHeight)

	rendered := rimg.NewRaster(imgWidth, imgHeight)

	var lastPercent float64 = 0.0
	renderStart := time.Now().UnixMilli()
	rtracer.TraceAllRays(scene, rendered, 100, 100, 0, func(n int32, percent float64) {
		if int(percent*100) >= int(lastPercent*100)+2 {
			lastPercent = percent
			ctx.RenderProgress.Set(percent)
		}
	})
	renderEnd := time.Now().UnixMilli()

	ctx.StatusText.Set(fmt.Sprintf("Done (%d ms).", renderEnd-renderStart))
	ctx.IsRendering = false
	ctx.RenderProgress.Set(0.0)
	ctx.PreviewImage.Update(rendered)
}
