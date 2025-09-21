package rmaterial

import (
	"encoding/json"

	"github.com/piotrwyrw/radia/radia/rcolor"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/rtypes"
)

type NormalMaterial struct {
	Brightness float64 `json:"brightness"`
}

func (e *NormalMaterial) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}

func NewNormalMaterial(brightness float64) *NormalMaterial {
	return &NormalMaterial{
		Brightness: brightness,
	}
}

func (m *NormalMaterial) Scatter(incoming *rmath.Ray, intersection *rtypes.Intersection) (scattered *rmath.Ray, attenuation rcolor.Color) {
	return nil, rcolor.ColorWhite()
}

func (e *NormalMaterial) Emitted(intersection *rtypes.Intersection) rcolor.Color {
	N := intersection.Object.Normal(intersection.Point)
	N.Normalize()

	C := rcolor.Color{
		R: (N.Y + 1) / 2.0,
		G: (N.Z + 1) / 2.0,
		B: (N.X + 1) / 2.0,
	}

	return C.MultiplyScalar(e.Brightness)
}

func (e *NormalMaterial) Identifier() string {
	return rtypes.NormalMaterialIdentifier
}
