package rmaterial

import (
	"github.com/piotrwyrw/radia/internal/rmath"
	"github.com/piotrwyrw/radia/radia/rcolor"
	"github.com/piotrwyrw/radia/radia/rscene"
)

type GlassMaterial struct {
	IOR float64
}

func (e *GlassMaterial) Scatter(incoming *rmath.Ray, intersection *rscene.Intersection) (scattered *rmath.Ray, attenuation rcolor.Color) {
	// TODO
	return &rmath.Ray{
		Origin:    intersection.Point.Copy(),
		Direction: intersection.Incoming.Direction,
	}, rcolor.ColorWhite()
}

func (e *GlassMaterial) Emitted(intersection *rscene.Intersection) rcolor.Color {
	return rcolor.ColorBlack()
}

func (e *GlassMaterial) Identifier() string {
	return "RadiaGlassMaterial"
}

func (e *GlassMaterial) Type() rscene.MaterialType {
	return rscene.ShapeMaterialType
}
