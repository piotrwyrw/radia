package rmaterial

import (
	"encoding/json"

	"github.com/piotrwyrw/radia/radia/rcolor"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/rtypes"
)

type GlassMaterial struct {
	IOR float64 `json:"ior"`
}

func (e *GlassMaterial) Unmarshal(data []byte) error {
	return json.Unmarshal(data, e)
}

func (e *GlassMaterial) Scatter(incoming *rmath.Ray, intersection *rtypes.Intersection) (scattered *rmath.Ray, attenuation rcolor.Color) {
	// TODO
	direction := incoming.Direction.Copy()
	direction.Y = -direction.Y
	direction.Normalize()
	return &rmath.Ray{
		Origin:    intersection.Point.Copy(),
		Direction: direction,
	}, rcolor.ColorWhite()
}

func (e *GlassMaterial) Emitted(intersection *rtypes.Intersection) rcolor.Color {
	return rcolor.ColorBlack()
}

func (e *GlassMaterial) Identifier() string {
	return rtypes.GlassMaterialIdentifier
}
