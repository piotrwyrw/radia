package rtypes

import "github.com/piotrwyrw/radia/radia/rmath"

type Camera struct {
	Location    rmath.Vec3d
	Facing      rmath.Vec3d
	FocalLength float64
}
