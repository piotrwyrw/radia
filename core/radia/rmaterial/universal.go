package rmaterial

import (
	"encoding/json"

	rmath2 "github.com/piotrwyrw/radia/internal/rmath"
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
	reflectDir := intersection.Object.Reflect(intersection)

	rand := rmath.RandomVector()
	rand.Multiply(m.Roughness)

	scatterDir := reflectDir.Copy()
	scatterDir.Add(rand)
	scatterDir.Normalize()

	scattered = &rmath.Ray{
		Origin:    intersection.Point.Copy(),
		Direction: scatterDir,
	}

	// Attenuation = base color scaled to prevent double-bright
	attenuation = m.Color.MultiplyScalar(rmath2.Clamp[float64, float64](0.0, 1.0, 1-m.Roughness))

	cosTheta := rmath2.Clamp[float64, float64](0.1, 1.0, scatterDir.Dot(intersection.Object.Normal(intersection.Point)))
	attenuation = attenuation.MultiplyScalar(cosTheta)

	return scattered, attenuation
}

func (e *UniversalMaterial) Emitted(intersection *rtypes.Intersection) rcolor.Color {
	emission := e.Emission
	emission = emission.MultiplyScalar(e.Brightness)
	return emission
}

func (e *UniversalMaterial) Identifier() string {
	return rtypes.UniversalMaterialIdentifier
}
