package scene

import (
	"raytracer/internal/rt/math"
)

type Intersection struct {
	Point    math.Vec3d
	Distance float64
	Incoming *math.Ray
	Object   Shape
}

type Shape interface {
	Hit(ray *math.Ray) *Intersection
	Normal(at math.Vec3d) math.Vec3d
	Reflect(intersection *Intersection) math.Vec3d
	GetMaterial() Material
}
