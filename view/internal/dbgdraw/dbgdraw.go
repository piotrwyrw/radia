package dbgdraw

import (
	"github.com/piotrwyrw/radia/radia/rimg"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/rtypes"
)
i
type DrawContext struct {
	Camera      *rtypes.Camera
	destination *rimg.Raster
}

func (ctx *DrawContext) drawLine(raster rimg.Raster, x, y float64) (screenX int32, screenY int32) {
}

func (ctx *DrawContext) ProjectPoint(point rmath.Vec3d) (screenX int64, screenY int64) {
	sX := (ctx.Camera.FocalLength * point.X) / point.Z
	sY := (ctx.Camera.FocalLength * point.Y) / point.Z
	return int64(sX), int64(sY)
}
