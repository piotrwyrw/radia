package rscene

import (
	rmath2 "github.com/piotrwyrw/radia/internal/rmath"
	"github.com/piotrwyrw/radia/radia/rmath"
)

type Shape interface {
	Hit(ray *rmath2.Ray) *Intersection
	Normal(at rmath.Vec3d) rmath.Vec3d
	Reflect(intersection *Intersection) rmath.Vec3d
	GetMaterial() ShapeMaterialWrapper
}

type Intersection struct {
	Point    rmath.Vec3d
	Distance float64
	Incoming *rmath2.Ray
	Object   Shape
}
