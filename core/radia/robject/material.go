package robject

import (
	"github.com/piotrwyrw/radia/radia/rtypes"
)

func WrapShapeMaterial(material rtypes.ShapeMaterial) rtypes.ShapeMaterialWrapper {
	return rtypes.ShapeMaterialWrapper{
		Type:     rtypes.ShapeWrapperType,
		Name:     material.Identifier(),
		Material: material,
	}
}

func WrapEnvironmentMaterial(material rtypes.EnvironmentMaterial) rtypes.EnvironmentMaterialWrapper {
	return rtypes.EnvironmentMaterialWrapper{
		Type:     rtypes.EnvironmentWrapperType,
		Name:     material.Identifier(),
		Material: material,
	}
}
