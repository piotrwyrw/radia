package scene

import "raytracer/internal/rt/math"

type Camera struct {
	Location    math.Vec3d
	Facing      math.Vec3d
	FocalLength float64
}
