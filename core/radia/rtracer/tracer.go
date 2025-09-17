package rtracer

import (
	math2 "math"

	rmath2 "github.com/piotrwyrw/radia/internal/rmath"
	"github.com/piotrwyrw/radia/radia/rcolor"
	"github.com/piotrwyrw/radia/radia/rimg"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/rscene"
)

func TraceAllRays(scene *rscene.Scene, raster *rimg.Raster, pixelSamples int, maxBounces int32, everyPixel func(x, y, n int32)) {
	// Find camera plane
	scene.Camera.Facing.Normalize()
	upDir := rmath.Vec3d{X: 0.0, Y: 1.0, Z: 0.0}

	rightDir := scene.Camera.Facing.Cross(upDir)
	rightDir.Normalize()

	trueUpDir := rightDir.Cross(scene.Camera.Facing)
	trueUpDir.Normalize()

	ar := float64(raster.Width) / float64(raster.Height)
	planeWidth := 1.0
	planeHeight := planeWidth / ar

	dx := planeWidth / float64(raster.Width)
	dy := planeHeight / float64(raster.Height)

	computeRayDirection := func(x, y float64) rmath.Vec3d {
		px := (-planeWidth / 2.0) + dx*x
		py := (planeHeight / 2.0) - dy*y

		dir := scene.Camera.Facing.Copy()
		dir.Multiply(scene.Camera.FocalLength)

		right := rightDir.Copy()
		right.Multiply(px)

		up := trueUpDir.Copy()
		up.Multiply(py)

		dir.Add(right)
		dir.Add(up)
		dir.Normalize()

		return dir
	}

	n := int32(0)

	// Trace all rays
	for x := int32(0); x < raster.Width; x++ {
		for y := int32(0); y < raster.Height; y++ {
			n++
			dir := computeRayDirection(float64(x), float64(y))
			colorSum := rcolor.Color{R: 0.0, G: 0.0, B: 0.0}
			for i := 0; i < pixelSamples; i++ {
				colorSum = colorSum.Add(TraceRay(scene.Camera.Location, dir, scene, 0, maxBounces))
			}
			raster.Put(x, y, colorSum.Divide(float64(pixelSamples)))
			go everyPixel(x, y, n)
		}
	}
}

func TraceRay(origin rmath.Vec3d, direction rmath.Vec3d, s *rscene.Scene, bounces int32, maxBounces int32) rcolor.Color {
	ray := rmath2.Ray{Origin: origin, Direction: direction}

	var distance = math2.Inf(1)
	var intersection *rscene.Intersection = nil

	for _, obj := range s.Objects {
		i := obj.Hit(&ray)
		if i == nil {
			continue
		}
		if i.Distance < distance {
			distance = i.Distance
			intersection = i
		}
	}

	if distance == math2.Inf(1) || intersection == nil {
		return s.WorldMat.Material.SkyColor(&direction)
	}

	mat := intersection.Object.GetMaterial().Material

	clr := mat.Emitted(intersection)
	scatter, attenuation := mat.Scatter(&ray, intersection)
	if scatter == nil || bounces > maxBounces {
		return clr
	}
	return clr.Add(attenuation.Multiply(TraceRay(intersection.Point, scatter.Direction, s, bounces+1, maxBounces)))
}
