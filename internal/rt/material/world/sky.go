package world

import (
	math2 "math"
	"raytracer/internal/rt/color"
	"raytracer/internal/rt/img"
	"raytracer/internal/rt/math"
)

type Sky struct {
	Image *img.Raster
	IOR   float64
}

func (sky *Sky) SkyColor(direction *math.Vec3d) color.Color {
	d := direction.Copy()
	d.Normalize()

	azimuth := math2.Atan2(d.Z, d.X) + math2.Pi*3
	for azimuth > math2.Pi {
		azimuth -= 2 * math2.Pi
	}
	for azimuth < -math2.Pi {
		azimuth += 2 * math2.Pi
	}

	elevation := math2.Acos(d.Y)

	u := (azimuth/(2*math2.Pi) + 0.5) * float64(sky.Image.Width)
	v := (elevation / math2.Pi) * float64(sky.Image.Height)

	uu := math.Clamp[float64, int32](0, float64(sky.Image.Width)-1, u)
	vv := math.Clamp[float64, int32](0, float64(sky.Image.Height)-1, v)

	px := sky.Image.Get(uu, vv)
	px = px.MultiplyScalar(2.0)
	return px
}
