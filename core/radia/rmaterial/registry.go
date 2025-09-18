package rmaterial

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/piotrwyrw/radia/radia/rtypes"
	"github.com/sirupsen/logrus"
)

type MaterialRegistry struct {
	environment map[string]reflect.Type
	shape       map[string]reflect.Type
}

var registryOnce sync.Once = sync.Once{}
var globalMatRegistry MaterialRegistry

func registerBuiltinMaterials() {
	_ = globalMatRegistry.Register(rtypes.UniversalMaterialIdentifier, reflect.TypeOf(&UniversalMaterial{}))
	_ = globalMatRegistry.Register(rtypes.GlassMaterialIdentifier, reflect.TypeOf(&GlassMaterial{}))
	_ = globalMatRegistry.Register(rtypes.SkyMaterialIdentifier, reflect.TypeOf(&Sky{}))
}

func GetMaterialRegistry() *MaterialRegistry {
	registryOnce.Do(func() {
		globalMatRegistry = newMaterialRegistry()
		registerBuiltinMaterials()
	})
	return &globalMatRegistry
}

func newMaterialRegistry() MaterialRegistry {
	return MaterialRegistry{
		environment: make(map[string]reflect.Type),
		shape:       make(map[string]reflect.Type),
	}
}

func (reg *MaterialRegistry) HasEnvironment(identifier string) bool {
	_, ok := reg.environment[identifier]
	return ok
}

func (reg *MaterialRegistry) HasShape(identifier string) bool {
	_, ok := reg.shape[identifier]
	return ok
}

func (reg *MaterialRegistry) InstantiateEnvironmentMaterial(identifier string) (rtypes.EnvironmentMaterial, error) {
	rt, ok := reg.environment[identifier]
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

func (reg *MaterialRegistry) InstantiateShapeMaterial(identifier string) (rtypes.ShapeMaterial, error) {
	rt, ok := reg.shape[identifier]
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

func (reg *MaterialRegistry) Register(identifier string, t reflect.Type) error {
	if t.Implements(reflect.TypeOf((*rtypes.EnvironmentMaterial)(nil)).Elem()) {
		reg.environment[identifier] = t
		logrus.Debugf("Registered environment material \"%s\".\n", identifier)
		return nil
	}
	if t.Implements(reflect.TypeOf((*rtypes.ShapeMaterial)(nil)).Elem()) {
		reg.shape[identifier] = t
		logrus.Debugf("Registered shape material \"%s\".\n", identifier)
		return nil
	}

	return fmt.Errorf("%s is not a material", identifier)
}
