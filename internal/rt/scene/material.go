package scene

import (
	"raytracer/internal/rt/color"
	"raytracer/internal/rt/math"
)

type Material interface {
	Scatter(incoming *math.Ray, intersection *Intersection) (scattered *math.Ray, attenuation color.Color)
	Emitted(intersection *Intersection) color.Color
}
