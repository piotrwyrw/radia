package rmaterial

import (
	"encoding/json"
	"math"

	"github.com/piotrwyrw/radia/radia/rcolor"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/rtypes"
)

type GradientSky struct {
	IOR          float64      `json:"ior"`
	Intensity    float64      `json:"intensity"`
	ColorHorizon rcolor.Color `json:"horizon_color"`
	ColorSky     rcolor.Color `json:"sky_color"`
}

func (sky *GradientSky) Unmarshal(data []byte) error {
	return json.Unmarshal(data, sky)
}

func (sky *GradientSky) SkyColor(direction *rmath.Vec3d) rcolor.Color {
	dir := direction.Copy()
	dir.Normalize()

	f := math.Acos(direction.Y) / math.Pi

	color := rcolor.ColorLerp(sky.ColorSky, sky.ColorHorizon, f/2)
	return color.MultiplyScalar(sky.Intensity)
}

func (sky *GradientSky) Identifier() string {
	return rtypes.GradientSkyMaterialIdentifier
}
