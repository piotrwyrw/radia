package rtypes

import (
	"github.com/piotrwyrw/radia/radia/rcolor"
	rmath2 "github.com/piotrwyrw/radia/radia/rmath"
)

const SkyMaterialIdentifier = "radia.sky"
const GradientSkyMaterialIdentifier = "radia.gradient_sky"
const NormalMaterialIdentifier = "radia.normal"
const GlassMaterialIdentifier = "radia.glass"
const UniversalMaterialIdentifier = "radia.universal"
const MirrorMaterialIdentifier = "radia.mirror"

const EnvironmentWrapperType = "environment"
const ShapeWrapperType = "shape"

type GenericMaterial interface {
	Identifier() string
}

type ShapeMaterial interface {
	Scatter(incoming *rmath2.Ray, intersection *Intersection) (scattered *rmath2.Ray, attenuation rcolor.Color)
	Emitted(intersection *Intersection) rcolor.Color
	Identifier() string
	Unmarshal(data []byte) error
}

type ShapeMaterialWrapper struct {
	Type     string        `json:"type"`
	Name     string        `json:"name"`
	Material ShapeMaterial `json:"material"`
}

type EnvironmentMaterial interface {
	SkyColor(direction *rmath2.Vec3d) rcolor.Color
	Identifier() string
	Unmarshal(data []byte) error
}

type EnvironmentMaterialWrapper struct {
	Type     string              `json:"type"`
	Name     string              `json:"name"`
	Material EnvironmentMaterial `json:"material"`
}
