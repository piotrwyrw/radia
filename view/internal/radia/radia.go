package radia

import (
	"github.com/piotrwyrw/otherproj/internal/context"
	"github.com/piotrwyrw/otherproj/internal/util"
	"github.com/piotrwyrw/radia/radia/radia"
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

	radia.Initialize()

	//mat := rmaterial.UniversalMaterial{
	//	Color:     rcolor.Color{R: 1.0, G: 0.5, B: 0.5},
	//	Emission:  rcolor.ColorBlack(),
	//	Roughness: 1.0,
	//}

	mat := rmaterial.GlassMaterial{
		IOR: 1.0,
	}

	scene := &rtypes.Scene{
		Objects: []rtypes.ShapeWrapper{
			robject.WrapShape(&rshapes.Sphere{
				Center:   rmath.Vec3d{X: 0.0, Y: 0.0, Z: 2.0},
				Radius:   0.3,
				Material: robject.WrapShapeMaterial(&mat),
			}),
		},
		Camera: rtypes.Camera{
			Location:    rmath.Vec3d{X: 0.0, Y: 0.0, Z: 0.0},
			Facing:      rmath.Vec3d{X: 0.0, Y: 0.0, Z: 1.0},
			FocalLength: 1.5,
		},
		WorldMat: robject.WrapEnvironmentMaterial(&rmaterial.Sky{
			Image:         env,
			IOR:           1.0,
			FallbackColor: rcolor.ColorBlack(),
		}),
	}

	_ = rscene.SaveSceneJSON(scene, "scene_js.json")
	_, err = rscene.LoadSceneJSON("scene_js.json")
	if err != nil {
		logrus.Errorf("Could not load scene json: %v\n", err)
	}

	destination := rimg.NewRaster(imageWidth, imageHeight)

	rtracer.TraceAllRays(scene, destination, 100, 100, func(x, y, n int32) {
		progress := float64(n) / float64(imageWidth*imageHeight)
		ctx.RenderProgress.Set(progress)
	})

	util.UpdateFyneImageWithRaster(destination, ctx.RenderOutputBuffer)

	ctx.StatusText.Set("Ready")
	ctx.RenderOutputImage.Refresh()

	ctx.IsRendering = false
	ctx.RenderProgress.Set(0.0)
}
