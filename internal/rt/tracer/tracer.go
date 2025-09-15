package tracer

import (
	math2 "math"
	"raytracer/internal/rt/color"
	"raytracer/internal/rt/img"
	"raytracer/internal/rt/math"
	"raytracer/internal/rt/scene"
)

func TraceAllRays(scene *scene.Scene, raster *img.Raster, pixelSamples int, maxBounces int32, everyPixel func(x, y, n int32)) {
	// Find camera plane
	scene.Camera.Facing.Normalize()
	upDir := math.Vec3d{X: 0.0, Y: 1.0, Z: 0.0}

	rightDir := scene.Camera.Facing.Cross(upDir)
	rightDir.Normalize()

	trueUpDir := rightDir.Cross(scene.Camera.Facing)
	trueUpDir.Normalize()

	ar := float64(raster.Width) / float64(raster.Height)
	planeWidth := 1.0
	planeHeight := planeWidth / ar

	dx := planeWidth / float64(raster.Width)
	dy := planeHeight / float64(raster.Height)

	computeRayDirection := func(x, y float64) math.Vec3d {
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
			colorSum := color.Color{R: 0.0, G: 0.0, B: 0.0}
			for i := 0; i < pixelSamples; i++ {
				colorSum = colorSum.Add(TraceRay(scene.Camera.Location, dir, scene, 0, maxBounces))
			}
			raster.Put(x, y, colorSum.Divide(float64(pixelSamples)))
			go everyPixel(x, y, n)
		}
	}
}

func TraceRay(origin math.Vec3d, direction math.Vec3d, s *scene.Scene, bounces int32, maxBounces int32) color.Color {
	ray := math.Ray{Origin: origin, Direction: direction}

	var distance = math2.Inf(1)
	var intersection *scene.Intersection = nil

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
		return s.WorldMat.SkyColor(&direction)
	}

	mat := intersection.Object.GetMaterial()

	clr := mat.Emitted(intersection)
	scatter, attenuation := mat.Scatter(&ray, intersection)
	if scatter == nil || bounces > maxBounces {
		return clr
	}
	return clr.Add(attenuation.Multiply(TraceRay(intersection.Point, scatter.Direction, s, bounces+1, maxBounces)))
}
