package rworld

import (
	math "math"

	irmath "github.com/piotrwyrw/radia/internal/rmath"
	"github.com/piotrwyrw/radia/radia/rcolor"
	"github.com/piotrwyrw/radia/radia/rimg"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/rscene"
)

type Sky struct {
	Image         *rimg.Raster `json:"image"`
	FallbackColor rcolor.Color `json:"fallbackColor"`
	IOR           float64      `json:"ior"`
}

func (sky *Sky) SkyColor(direction *rmath.Vec3d) rcolor.Color {
	if sky.Image == nil {
		return sky.FallbackColor
	}

	d := direction.Copy()
	d.Normalize()

	azimuth := math.Atan2(d.Z, d.X) + math.Pi*3
	for azimuth > math.Pi {
		azimuth -= 2 * math.Pi
	}
	for azimuth < -math.Pi {
		azimuth += 2 * math.Pi
	}

	elevation := math.Acos(d.Y)

	u := (azimuth/(2*math.Pi) + 0.5) * float64(sky.Image.Width)
	v := (elevation / math.Pi) * float64(sky.Image.Height)

	uu := irmath.Clamp[float64, int32](0, float64(sky.Image.Width)-1, u)
	vv := irmath.Clamp[float64, int32](0, float64(sky.Image.Height)-1, v)

	px := sky.Image.Get(uu, vv)
	px = px.MultiplyScalar(2.0)
	return px
}

func (sky *Sky) Identifier() string {
	return "RadiaSky"
}

func (sky *Sky) Type() rscene.MaterialType {
	return rscene.EnvironmentMaterialType
}
