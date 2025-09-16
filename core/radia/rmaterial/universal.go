package rmaterial

import (
	rmath2 "github.com/piotrwyrw/radia/internal/rmath"
	"github.com/piotrwyrw/radia/radia/rcolor"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/rscene"
)

type UniversalMaterial struct {
	Color            rcolor.Color
	Emission         rcolor.Color
	EmissionStrength float64
	Roughness        float64
}

func NewUniversalMaterial(color rcolor.Color, emission rcolor.Color, emissionStrength float64, roughness float64) *UniversalMaterial {
	return &UniversalMaterial{
		Color:            color,
		Emission:         emission,
		EmissionStrength: emissionStrength,
		Roughness:        roughness,
	}
}

func (e *UniversalMaterial) Scatter(incoming *rmath2.Ray, intersection *rscene.Intersection) (scattered *rmath2.Ray, attenuation rcolor.Color) {
	reflected := intersection.Object.Reflect(intersection)

	// Random direction in hemisphere
	scatterDir := reflected.Copy()
	rand := rmath.RandomVector()
	rand.Multiply(e.Roughness)
	scatterDir.Add(rand)
	scatterDir.Normalize()

	return &rmath2.Ray{
		Origin:    intersection.Point.Copy(),
		Direction: scatterDir,
	}, e.Color
}

func (e *UniversalMaterial) Emitted(intersection *rscene.Intersection) rcolor.Color {
	emission := e.Emission
	emission = emission.MultiplyScalar(e.EmissionStrength)
	return emission
}

func (e *UniversalMaterial) Identifier() string {
	return "RadiaUniversalMaterial"
}

func (e *UniversalMaterial) Type() rscene.MaterialType {
	return rscene.ShapeMaterialType
}
