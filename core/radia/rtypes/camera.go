package rtypes

import "github.com/piotrwyrw/radia/radia/rmath"

type Camera struct {
	Location    rmath.Vec3d `json:"location"`
	Facing      rmath.Vec3d `json:"facing"`
	FocalLength float64     `json:"focal_length" ui:"Focal Length"`
}
