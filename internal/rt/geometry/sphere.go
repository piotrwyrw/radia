package geometry

import (
	math2 "math"
	"raytracer/internal/rt/math"
	"raytracer/internal/rt/scene"
)

type Sphere struct {
	Center   math.Vec3d
	Radius   float64
	Material scene.Material
}

func (s *Sphere) Hit(ray *math.Ray) *scene.Intersection {
	L := ray.Origin.Copy()
	L.Sub(s.Center)

	a := ray.Direction.Dot(ray.Direction)
	b := 2 * (L.Dot(ray.Direction))
	c := L.Dot(L) - s.Radius*s.Radius

	solutions := math.Quadratic(a, b, c)

	if len(solutions) == 0 {
		return nil
	}

	var distance = math2.Inf(1)
	for _, t := range solutions {
		if t > 0 && t < distance {
			distance = t
		}
	}
	if distance <= math.EPSILON || distance == math2.Inf(1) {
		return nil
	}

	point := ray.Direction.Copy()
	point.Resize(distance)
	point.Add(ray.Origin)

	return &scene.Intersection{
		Point:    point,
		Distance: distance,
		Object:   s,
		Incoming: ray,
	}
}

func (s *Sphere) Normal(at math.Vec3d) math.Vec3d {
	normal := at.Copy()
	normal.Sub(s.Center)
	normal.Normalize()
	return normal
}

func (s *Sphere) Reflect(intersection *scene.Intersection) math.Vec3d {
	normal := s.Normal(intersection.Point)
	incoming := intersection.Incoming.Direction

	vDotN := incoming.Dot(normal)
	scaledNormal := normal.Copy()
	scaledNormal.Multiply(2 * vDotN)

	reflected := incoming.Copy()
	reflected.Sub(scaledNormal)

	return reflected
}

func (s *Sphere) GetMaterial() scene.Material {
	return s.Material
}
