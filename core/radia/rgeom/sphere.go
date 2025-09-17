package rgeom

import (
	math2 "math"

	rmath2 "github.com/piotrwyrw/radia/internal/rmath"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/rscene"
)

type Sphere struct {
	Center   rmath.Vec3d
	Radius   float64
	Material rscene.ShapeMaterialWrapper
}

func (s *Sphere) Hit(ray *rmath2.Ray) *rscene.Intersection {
	L := ray.Origin.Copy()
	L.Sub(s.Center)

	a := ray.Direction.Dot(ray.Direction)
	b := 2 * (L.Dot(ray.Direction))
	c := L.Dot(L) - s.Radius*s.Radius

	solutions := rmath2.Quadratic(a, b, c)

	if len(solutions) == 0 {
		return nil
	}

	var distance = math2.Inf(1)
	for _, t := range solutions {
		if t > 0 && t < distance {
			distance = t
		}
	}
	if distance <= rmath2.EPSILON || distance == math2.Inf(1) {
		return nil
	}

	point := ray.Direction.Copy()
	point.Resize(distance)
	point.Add(ray.Origin)

	return &rscene.Intersection{
		Point:    point,
		Distance: distance,
		Object:   s,
		Incoming: ray,
	}
}

func (s *Sphere) Normal(at rmath.Vec3d) rmath.Vec3d {
	normal := at.Copy()
	normal.Sub(s.Center)
	normal.Normalize()
	return normal
}

func (s *Sphere) Reflect(intersection *rscene.Intersection) rmath.Vec3d {
	normal := s.Normal(intersection.Point)
	incoming := intersection.Incoming.Direction

	vDotN := incoming.Dot(normal)
	scaledNormal := normal.Copy()
	scaledNormal.Multiply(2 * vDotN)

	reflected := incoming.Copy()
	reflected.Sub(scaledNormal)

	return reflected
}

func (s *Sphere) GetMaterial() rscene.ShapeMaterialWrapper {
	return s.Material
}
