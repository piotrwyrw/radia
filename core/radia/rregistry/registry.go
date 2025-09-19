package rregistry

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/piotrwyrw/radia/radia/rmaterial"
	"github.com/piotrwyrw/radia/radia/rshapes"
	"github.com/piotrwyrw/radia/radia/rtypes"
	"github.com/sirupsen/logrus"
)

type CentralRegistry struct {
	envMat   map[string]reflect.Type
	shapeMat map[string]reflect.Type
	shapes   map[string]reflect.Type
}

var registryOnce sync.Once = sync.Once{}
var globalMatRegistry CentralRegistry

func registerBuiltinMaterials() {
	// Shape Material
	_ = globalMatRegistry.Register(rtypes.UniversalMaterialIdentifier, reflect.TypeOf(&rmaterial.UniversalMaterial{}))
	_ = globalMatRegistry.Register(rtypes.GlassMaterialIdentifier, reflect.TypeOf(&rmaterial.GlassMaterial{}))

	// Environment Material
	_ = globalMatRegistry.Register(rtypes.SkyMaterialIdentifier, reflect.TypeOf(&rmaterial.Sky{}))
	_ = globalMatRegistry.Register(rtypes.GradientSkyMaterialIdentifier, reflect.TypeOf(&rmaterial.GradientSky{}))

	// Shapes
	_ = globalMatRegistry.Register(rtypes.ShapeIdentifierSphere, reflect.TypeOf(&rshapes.Sphere{}))
}

func GetCentralRegistry() *CentralRegistry {
	registryOnce.Do(func() {
		globalMatRegistry = newCentralRegistry()
		registerBuiltinMaterials()
	})
	return &globalMatRegistry
}

func newCentralRegistry() CentralRegistry {
	return CentralRegistry{
		envMat:   make(map[string]reflect.Type),
		shapeMat: make(map[string]reflect.Type),
		shapes:   make(map[string]reflect.Type),
	}
}

func (reg *CentralRegistry) HasEnvMat(identifier string) bool {
	_, ok := reg.envMat[identifier]
	return ok
}

func (reg *CentralRegistry) HasShapeMat(identifier string) bool {
	_, ok := reg.shapeMat[identifier]
	return ok
}

func (reg *CentralRegistry) HasShape(identifier string) bool {
	_, ok := reg.shapes[identifier]
	return ok
}

func (reg *CentralRegistry) InstantiateEnvMat(identifier string) (rtypes.EnvironmentMaterial, error) {
	rt, ok := reg.envMat[identifier]
	if !ok {
		return nil, fmt.Errorf("no environment material found for identifier: %s", identifier)
	}
	var mat rtypes.EnvironmentMaterial
	if rt.Kind() == reflect.Ptr {
		v := reflect.New(rt.Elem()).Interface()
		mat = v.(rtypes.EnvironmentMaterial)
	} else {
		v := reflect.New(rt).Interface()
		mat = v.(rtypes.EnvironmentMaterial)
	}
	return mat, nil
}

func (reg *CentralRegistry) InstantiateShapeMat(identifier string) (rtypes.ShapeMaterial, error) {
	rt, ok := reg.shapeMat[identifier]
	if !ok {
		return nil, fmt.Errorf("no shape material found for identifier: %s", identifier)
	}
	var mat rtypes.ShapeMaterial
	if rt.Kind() == reflect.Ptr {
		v := reflect.New(rt.Elem()).Interface()
		mat = v.(rtypes.ShapeMaterial)
	} else {
		v := reflect.New(rt).Interface()
		mat = v.(rtypes.ShapeMaterial)
	}
	return mat, nil
}

func (reg *CentralRegistry) InstantiateShape(identifier string) (rtypes.Shape, error) {
	st, ok := reg.shapes[identifier]
	if !ok {
		return nil, fmt.Errorf("no shape found for identifier: %s", identifier)
	}
	var s rtypes.Shape
	if st.Kind() == reflect.Ptr {
		v := reflect.New(st.Elem()).Interface()
		s = v.(rtypes.Shape)
	} else {
		v := reflect.New(st).Interface()
		s = v.(rtypes.Shape)
	}
	return s, nil
}

func (reg *CentralRegistry) Register(identifier string, t reflect.Type) error {
	if t.Implements(reflect.TypeOf((*rtypes.EnvironmentMaterial)(nil)).Elem()) {
		reg.envMat[identifier] = t
		logrus.Debugf("Registered environment material: \"%s\"\n", identifier)
		return nil
	}

	if t.Implements(reflect.TypeOf((*rtypes.ShapeMaterial)(nil)).Elem()) {
		reg.shapeMat[identifier] = t
		logrus.Debugf("Registered shape material: \"%s\"\n", identifier)
		return nil
	}

	if t.Implements(reflect.TypeOf((*rtypes.Shape)(nil)).Elem()) {
		reg.shapes[identifier] = t
		logrus.Debugf("Registered shape: \"%s\"\n", identifier)
		return nil
	}

	return fmt.Errorf("\"%s\" (Type \"%s\") is neither a material, nor a shape", identifier, t.Name())
}
