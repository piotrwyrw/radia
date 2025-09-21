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

func (m *UniversalMaterial) Scatter(incoming *rmath.Ray, intersection *rtypes.Intersection) (scattered *rmath.Ray, attenuation rcolor.Color) {
	N := intersection.Object.Normal(intersection.Point)

	rand := rmath.RandomVector()
	rand.Multiply(m.Roughness)

	diffuseDir := N.Copy()
	diffuseDir.Add(rand)
	diffuseDir.Normalize()

	scattered = &rmath.Ray{
		Origin:    intersection.Point.Copy(),
		Direction: diffuseDir,
	}

	att := m.Color
	att.MultiplyScalar(0.8)

	return scattered, att
}

func (e *UniversalMaterial) Emitted(intersection *rtypes.Intersection) rcolor.Color {
	emission := e.Emission
	emission = emission.MultiplyScalar(e.Brightness)
	return emission
}

func (e *UniversalMaterial) Identifier() string {
	return rtypes.UniversalMaterialIdentifier
}
