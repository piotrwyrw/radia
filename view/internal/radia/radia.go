package radia

import (
	"fmt"

	"github.com/piotrwyrw/otherproj/internal/context"
	"github.com/piotrwyrw/otherproj/internal/util"
	"github.com/piotrwyrw/radia/radia/rcolor"
	"github.com/piotrwyrw/radia/radia/rgeom"
	"github.com/piotrwyrw/radia/radia/rimg"
	"github.com/piotrwyrw/radia/radia/rmaterial"
	"github.com/piotrwyrw/radia/radia/rmaterial/rworld"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/rscene"
	"github.com/piotrwyrw/radia/radia/rtracer"
)

func InvokeRenderer(ctx *context.Context, imageWidth int32, imageHeight int32) {
	ctx.StatusText.Set("Rendering ...")
	ctx.IsRendering = true
	ctx.RenderProgress.Set(0.0)

	env, err := rimg.RasterFromPNG("world.png")
	if err != nil {
		fmt.Printf("Could not render: %v\n", err)
		ctx.IsRendering = false
		return
	}

	mat := rmaterial.UniversalMaterial{
		Color:     rcolor.Color{R: 1.0, G: 0.5, B: 0.5},
		Emission:  rcolor.ColorBlack(),
		Roughness: 0.5,
	}

	scene := rscene.Scene{
		Objects: []rscene.Shape{
			&rgeom.Sphere{
				Center:   rmath.Vec3d{X: 0.0, Y: 0.0, Z: 2.0},
				Radius:   0.3,
				Material: rscene.WrapShapeMaterial(&mat),
			},
		},
		Camera: rscene.Camera{
			Location:    rmath.Vec3d{X: 0.0, Y: 0.0, Z: 0.0},
			Facing:      rmath.Vec3d{X: 0.0, Y: 0.0, Z: 1.0},
			FocalLength: 1.5,
		},
		WorldMat: rscene.WrapEnvironmentMaterial(&rworld.Sky{
			Image:         env,
			IOR:           1.0,
			FallbackColor: rcolor.ColorBlack(),
		}),
	}

	_ = scene.SaveJSON("scene_js.json")

	destination := rimg.NewRaster(imageWidth, imageHeight)

	rtracer.TraceAllRays(&scene, destination, 10, 10, func(x, y, n int32) {
		progress := float64(n) / float64(imageWidth*imageHeight)
		ctx.RenderProgress.Set(progress)
	})

	util.UpdateFyneImageWithRaster(destination, ctx.RenderOutputBuffer)

	ctx.StatusText.Set("Ready")
	ctx.RenderOutputImage.Refresh()

	ctx.IsRendering = false
	ctx.RenderProgress.Set(0.0)
}
