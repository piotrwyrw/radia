package rmaterial

import (
	"encoding/json"
	math "math"

	irmath "github.com/piotrwyrw/radia/internal/rmath"
	"github.com/piotrwyrw/radia/radia/rcolor"
	"github.com/piotrwyrw/radia/radia/rimg"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/rtypes"
)

type Sky struct {
	Image         *rimg.Raster `json:"image"`
	FallbackColor rcolor.Color `json:"fallback_color"`
	IOR           float64      `json:"ior"`
	Azimuth       float64      `json:"azimuth"`
	Elevation     float64      `json:"elevation"`
	Intensity     float64      `json:"intensity"`
}

func (sky *Sky) Unmarshal(data []byte) error {
	return json.Unmarshal(data, sky)
}

func (sky *Sky) SkyColor(direction *rmath.Vec3d) rcolor.Color {
	if sky.Image == nil {
		return sky.FallbackColor
	}

	d := direction.Copy()
	d.Normalize()

	azimuth := math.Atan2(d.Z, d.X) + sky.Azimuth
	for azimuth > math.Pi {
		azimuth -= 2 * math.Pi
	}
	for azimuth < -math.Pi {
		azimuth += 2 * math.Pi
	}

	elevation := math.Acos(d.Y) + sky.Elevation

	u := (0.5 + azimuth/(math.Pi*2)) * float64(sky.Image.Width)
	v := (elevation / math.Pi) * float64(sky.Image.Height)

	uu := irmath.Clamp[float64, int32](0, float64(sky.Image.Width)-1, u)
	vv := irmath.Clamp[float64, int32](0, float64(sky.Image.Height)-1, v)

	px := sky.Image.Get(uu, vv)
	return px.MultiplyScalar(sky.Intensity)
}

func (sky *Sky) Identifier() string {
	return rtypes.SkyMaterialIdentifier
}
