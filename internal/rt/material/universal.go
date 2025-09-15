package material

import (
	"raytracer/internal/rt/color"
	"raytracer/internal/rt/math"
	"raytracer/internal/rt/scene"
)

type UniversalMaterial struct {
	Color            color.Color
	Emission         color.Color
	EmissionStrength float64
	Roughness        float64
}

func NewUniversalMaterial(color color.Color, emission color.Color, emissionStrength float64, roughness float64) *UniversalMaterial {
	return &UniversalMaterial{
		Color:            color,
		Emission:         emission,
		EmissionStrength: emissionStrength,
		Roughness:        roughness,
	}
}

func (e *UniversalMaterial) Scatter(incoming *math.Ray, intersection *scene.Intersection) (scattered *math.Ray, attenuation color.Color) {
	reflected := intersection.Object.Reflect(intersection)

	// Random direction in hemisphere
	scatterDir := reflected.Copy()
	rand := math.RandomVector()
	rand.Multiply(e.Roughness)
	scatterDir.Add(rand)
	scatterDir.Normalize()

	return &math.Ray{
		Origin:    intersection.Point.Copy(),
		Direction: scatterDir,
	}, e.Color
}

func (e *UniversalMaterial) Emitted(intersection *scene.Intersection) color.Color {
	emission := e.Emission
	emission = emission.MultiplyScalar(e.EmissionStrength)
	return emission
}
