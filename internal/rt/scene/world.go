package scene

import (
	"raytracer/internal/rt/color"
	"raytracer/internal/rt/math"
)

type WorldMaterial interface {
	SkyColor(direction *math.Vec3d) color.Color
}
