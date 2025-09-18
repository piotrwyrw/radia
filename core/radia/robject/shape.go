package robject

import (
	"github.com/piotrwyrw/radia/radia/rtypes"
)

func WrapShape(shape rtypes.Shape) rtypes.ShapeWrapper {
	return rtypes.ShapeWrapper{
		Type:   shape.Identifier(),
		Object: shape.(rtypes.Shape),
	}
}
