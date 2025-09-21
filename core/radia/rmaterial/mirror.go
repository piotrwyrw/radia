package rmaterial

import (
	"encoding/json"

	"github.com/piotrwyrw/radia/radia/rcolor"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/rtypes"
)

type MirrorMaterial struct {
	Color rcolor.Color
}

func (e *MirrorMaterial) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}

func NewMirrorMaterial() *MirrorMaterial {
	return &MirrorMaterial{}
}

func (m *MirrorMaterial) Scatter(incoming *rmath.Ray, intersection *rtypes.Intersection) (scattered *rmath.Ray, attenuation rcolor.Color) {
	R := intersection.Object.Reflect(intersection)
	R.Normalize()
	return &rmath.Ray{
		Origin:    intersection.Point,
		Direction: R,
	}, m.Color
}

func (e *MirrorMaterial) Emitted(intersection *rtypes.Intersection) rcolor.Color {
	return rcolor.ColorBlack()
}

func (e *MirrorMaterial) Identifier() string {
	return rtypes.MirrorMaterialIdentifier
}
