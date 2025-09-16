package rscene

import (
	"github.com/piotrwyrw/radia/internal/rmath"
	"github.com/piotrwyrw/radia/radia/rcolor"
	rmath2 "github.com/piotrwyrw/radia/radia/rmath"
)

type MaterialType int

const (
	ShapeMaterialType = iota
	EnvironmentMaterialType
)

type GenericMaterial interface {
	Identifier() string
	Type() MaterialType
}

type ShapeMaterial interface {
	Scatter(incoming *rmath.Ray, intersection *Intersection) (scattered *rmath.Ray, attenuation rcolor.Color)
	Emitted(intersection *Intersection) rcolor.Color

	Identifier() string
	Type() MaterialType
}

type EnvironmentMaterial interface {
	SkyColor(direction *rmath2.Vec3d) rcolor.Color

	Identifier() string
	Type() MaterialType
}
