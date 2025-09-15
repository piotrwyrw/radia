package main

import (
	"fmt"
	math2 "math"
	"math/rand"
	"os"
	"raytracer/internal/rt/color"
	"raytracer/internal/rt/geometry"
	"raytracer/internal/rt/img"
	"raytracer/internal/rt/material"
	"raytracer/internal/rt/material/world"
	"raytracer/internal/rt/math"
	"raytracer/internal/rt/scene"
	"raytracer/internal/viewer"
)

func main() {
	env, err := img.RasterFromPNG("world.png")
	if err != nil {
		fmt.Printf("Could not load environemnt texture: %s\n", err)
		os.Exit(1)
		return
	}

	var objects []scene.Shape = []scene.Shape{
		// Ground plane
		&geometry.Sphere{
			Center:   math.Vec3d{X: 0.0, Y: -1000.0, Z: 0.0},
			Radius:   1000.0,
			Material: material.NewUniversalMaterial(color.ColorWhite(), color.ColorBlack(), 0.0, 2.0),
		},
	}

	// Generate random balls
	for i := 0; i < 150; i++ {
		angle := rand.Float64() * 2 * math2.Pi
		distance := rand.Float64()*2 + 1

		radius := rand.Float64()*0.1 + 0.05
		y := radius

		var mat scene.Material
		if rand.Int31n(10) > 5 {
			mat = material.NewUniversalMaterial(color.ColorRandom(), color.ColorBlack(), 0.0, rand.Float64())
		} else {
			mat = &material.GlassMaterial{IOR: 1.5}
		}

		objects = append(objects, &geometry.Sphere{
			Center:   math.Vec3d{X: math2.Cos(angle) * distance, Y: y, Z: math2.Sin(angle) * distance},
			Radius:   radius,
			Material: mat,
		})
	}

	s := scene.Scene{
		Objects: objects,
		Camera: scene.Camera{
			Location:    math.Vec3d{X: 0.0, Y: 0.1, Z: 0.0},
			Facing:      math.Vec3d{X: 0.0, Y: 0.0, Z: 1.0},
			FocalLength: 1.0,
		},
		WorldMat: &world.Sky{Image: env, IOR: 1.0},
	}
	err = viewer.StartRaytracer(1500, 900, &s)
	if err != nil {
		fmt.Printf("Failed: %s\n", err)
		os.Exit(1)
	}
}
