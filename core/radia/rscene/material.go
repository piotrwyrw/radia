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
}

type ShapeMaterial interface {
	Scatter(incoming *rmath.Ray, intersection *Intersection) (scattered *rmath.Ray, attenuation rcolor.Color)
	Emitted(intersection *Intersection) rcolor.Color
	Identifier() string
}

type ShapeMaterialWrapper struct {
	Identifier string
	Material   ShapeMaterial
}

func WrapShapeMaterial(material ShapeMaterial) ShapeMaterialWrapper {
	return ShapeMaterialWrapper{
		Identifier: material.Identifier(),
		Material:   material,
	}
}

type EnvironmentMaterial interface {
	SkyColor(direction *rmath2.Vec3d) rcolor.Color
	Identifier() string
}

type EnvironmentMaterialWrapper struct {
	Identifier string
	Material   EnvironmentMaterial
}

func WrapEnvironmentMaterial(material EnvironmentMaterial) EnvironmentMaterialWrapper {
	return EnvironmentMaterialWrapper{
		Identifier: material.Identifier(),
		Material:   material,
	}
}
