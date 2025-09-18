package rmaterial

import (
	"encoding/json"

	"github.com/piotrwyrw/radia/radia/rcolor"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/rtypes"
)

type UniversalMaterial struct {
	Color      rcolor.Color `json:"color"`
	Emission   rcolor.Color `json:"emission"`
	Brightness float64      `json:"brightness"`
	Roughness  float64      `json:"roughness"`
}

func (e *UniversalMaterial) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}

func NewUniversalMaterial(color rcolor.Color, emission rcolor.Color, brightness float64, roughness float64) *UniversalMaterial {
	return &UniversalMaterial{
		Color:      color,
		Emission:   emission,
		Brightness: brightness,
		Roughness:  roughness,
	}
}

func (e *UniversalMaterial) Scatter(incoming *rmath.Ray, intersection *rtypes.Intersection) (scattered *rmath.Ray, attenuation rcolor.Color) {
	reflected := intersection.Object.Reflect(intersection)

	// Random direction in hemisphere
	scatterDir := reflected.Copy()
	rand := rmath.RandomVector()
	rand.Multiply(e.Roughness)
	scatterDir.Add(rand)
	scatterDir.Normalize()

	return &rmath.Ray{
		Origin:    intersection.Point.Copy(),
		Direction: scatterDir,
	}, e.Color
}

func (e *UniversalMaterial) Emitted(intersection *rtypes.Intersection) rcolor.Color {
	emission := e.Emission
	emission = emission.MultiplyScalar(e.Brightness)
	return emission
}

func (e *UniversalMaterial) Identifier() string {
	return rtypes.UniversalMaterialIdentifier
}
