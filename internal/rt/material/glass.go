package material

import (
	"raytracer/internal/rt/color"
	"raytracer/internal/rt/math"
	"raytracer/internal/rt/scene"
)

type GlassMaterial struct {
	IOR float64
}

func (e *GlassMaterial) Scatter(incoming *math.Ray, intersection *scene.Intersection) (scattered *math.Ray, attenuation color.Color) {
	// TODO
	return &math.Ray{
		Origin:    intersection.Point.Copy(),
		Direction: intersection.Incoming.Direction,
	}, color.ColorWhite()
}

func (e *GlassMaterial) Emitted(intersection *scene.Intersection) color.Color {
	return color.ColorBlack()
}
