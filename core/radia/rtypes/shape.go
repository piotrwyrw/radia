package rtypes

import (
	"github.com/piotrwyrw/radia/radia/rmath"
)

const ShapeIdentifierSphere = "sphere"

type Shape interface {
	Hit(ray *rmath.Ray) *Intersection
	Normal(at rmath.Vec3d) rmath.Vec3d
	Reflect(intersection *Intersection) rmath.Vec3d
	GetMaterial() int32
	Identifier() string
	Unmarshal(data []byte) error
}

type ShapeWrapper struct {
	Type   string `json:"type"`
	Object Shape  `json:"object"`
}

type Intersection struct {
	Point    rmath.Vec3d
	Distance float64
	Incoming *rmath.Ray
	Object   Shape
}
