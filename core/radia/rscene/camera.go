package rscene

import (
	math2 "github.com/piotrwyrw/radia/radia/rmath"
)

type Camera struct {
	Location    math2.Vec3d
	Facing      math2.Vec3d
	FocalLength float64
}
