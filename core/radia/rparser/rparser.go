package rparser

import (
	"encoding/json"
	"fmt"

	"github.com/piotrwyrw/radia/radia/rmaterial"
	"github.com/piotrwyrw/radia/radia/robject"
	"github.com/piotrwyrw/radia/radia/rshapes"
	"github.com/piotrwyrw/radia/radia/rtypes"
)

func parseShapeWrapper(data []byte, dst *rtypes.ShapeWrapper, registry *rmaterial.MaterialRegistry) error {
	var aux struct {
		Type   string          `json:"type"`
		Object json.RawMessage `json:"object"`
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	dst.Type = aux.Type

	switch aux.Type {
	case rtypes.ShapeIdentifierSphere:
		var sphere rshapes.Sphere
		err = sphere.Unmarshal(aux.Object, func(data []byte, dst *rtypes.ShapeMaterialWrapper) error {
			return parseShapeMaterialWrapper(data, dst, registry)
		})
		dst.Object = &sphere
		break
	default:
		return fmt.Errorf("unknown shape type: %s", aux.Type)
	}

	if err != nil {
		return err
	}

	return nil
}

func parseShapeMaterialWrapper(data []byte, dst *rtypes.ShapeMaterialWrapper, registry *rmaterial.MaterialRegistry) error {
	var aux struct {
		Type     string          `json:"type"`
		Name     string          `json:"name"`
		Material json.RawMessage `json:"material"`
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	if aux.Type != rtypes.ShapeWrapperType {
		return fmt.Errorf("expected to parse shape material, got %s", aux.Type)
	}

	if !registry.HasShape(aux.Name) {
		return fmt.Errorf("shape material \"%s\" not found in registry", aux.Name)
	}

	mat, err := registry.InstantiateShapeMaterial(aux.Name)
	if err != nil {
		return err
	}
	err = mat.Unmarshal(aux.Material)
	if err != nil {
		return err
	}

	*dst = robject.WrapShapeMaterial(mat)

	return nil
}

func parseEnvironmentMaterialWrapper(data []byte, dst *rtypes.EnvironmentMaterialWrapper, registry *rmaterial.MaterialRegistry) error {
	var aux struct {
		Type     string          `json:"type"`
		Name     string          `json:"name"`
		Material json.RawMessage `json:"material"`
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	if aux.Type != rtypes.EnvironmentWrapperType {
		return fmt.Errorf("expected to parse environment material, got %s", aux.Type)
	}

	if !registry.HasEnvironment(aux.Name) {
		return fmt.Errorf("environment material \"%s\" not found in registry", aux.Name)
	}

	mat, err := registry.InstantiateEnvironmentMaterial(aux.Name)

	if err != nil {
		return err
	}
	err = mat.Unmarshal(aux.Material)
	if err != nil {
		return err
	}

	*dst = robject.WrapEnvironmentMaterial(mat)

	return nil
}

func ParseScene(data []byte, registry *rmaterial.MaterialRegistry) (*rtypes.Scene, error) {
	var aux struct {
		Metadata rtypes.SceneMetadata `json:"metadata"`
		Objects  []json.RawMessage    `json:"objects"`
		Camera   rtypes.Camera        `json:"camera"`
		WorldMat json.RawMessage      `json:"world"`
	}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return nil, err
	}

	var scene rtypes.Scene
	scene.Metadata = aux.Metadata
	scene.Camera = aux.Camera

	// Parse world material
	err = parseEnvironmentMaterialWrapper(aux.WorldMat, &scene.WorldMat, registry)
	if err != nil {
		return nil, err
	}

	// Parse scene objects
	for _, object := range aux.Objects {
		var shape rtypes.ShapeWrapper
		err = parseShapeWrapper(object, &shape, registry)
		if err != nil {
			return nil, err
		}
		scene.Objects = append(scene.Objects, shape)
	}

	return &scene, nil
}
