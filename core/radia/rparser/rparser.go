package rparser

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/piotrwyrw/radia/radia/robject"
	"github.com/piotrwyrw/radia/radia/rregistry"
	"github.com/piotrwyrw/radia/radia/rscene"
	"github.com/piotrwyrw/radia/radia/rtypes"
	"github.com/sirupsen/logrus"
)

func parseShapeWrapper(data []byte, dst *rtypes.ShapeWrapper, registry *rregistry.CentralRegistry) error {
	var aux struct {
		Type   string          `json:"type"`
		Object json.RawMessage `json:"object"`
	}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	if !registry.HasShape(aux.Type) {
		return fmt.Errorf("unknown shape: \"%s\"", aux.Type)
	}

	dst.Type = aux.Type

	s, err := registry.InstantiateShape(aux.Type)
	if err != nil {
		return err
	}
	err = s.Unmarshal(aux.Object)

	if err != nil {
		return err
	}

	dst.Object = s

	return nil
}

func parseShapeMaterialWrapper(data []byte, dst *rtypes.ShapeMaterialWrapper, registry *rregistry.CentralRegistry) error {
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

	if !registry.HasShapeMat(aux.Name) {
		return fmt.Errorf("unknown surface material: \"%s\"", aux.Name)
	}

	mat, err := registry.InstantiateShapeMat(aux.Name)
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

func parseEnvironmentMaterialWrapper(data []byte, dst *rtypes.EnvironmentMaterialWrapper, registry *rregistry.CentralRegistry) error {
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

	if !registry.HasEnvMat(aux.Name) {
		return fmt.Errorf("unknown environment material: \"%s\"", aux.Name)
	}

	mat, err := registry.InstantiateEnvMat(aux.Name)

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

func ParseScene(data []byte, registry *rregistry.CentralRegistry) (*rtypes.Scene, error) {
	var aux struct {
		Metadata  rtypes.SceneMetadata      `json:"metadata"`
		Materials map[int32]json.RawMessage `json:"materials"`
		Objects   []json.RawMessage         `json:"objects"`
		Camera    rtypes.Camera             `json:"camera"`
		WorldMat  json.RawMessage           `json:"world"`
	}
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return nil, err
	}

	var scene rtypes.Scene
	scene.Metadata = aux.Metadata
	scene.Camera = aux.Camera
	scene.Materials = make(map[int32]rtypes.ShapeMaterialWrapper, len(aux.Materials))

	// Parse materials
	for id, mat := range aux.Materials {
		var parsed rtypes.ShapeMaterialWrapper
		err := parseShapeMaterialWrapper(mat, &parsed, registry)
		if err != nil {
			return nil, err
		}

		// This will probably never happen. Unmarshalling doesn't preserve duplicate map entries
		if rscene.SceneMaterialExists(&scene, id) {
			return nil, fmt.Errorf("duplicate shape material: %d", id)
		}

		scene.Materials[id] = parsed
	}

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
		matId := shape.Object.GetMaterial()
		if !rscene.SceneMaterialExists(&scene, matId) {
			return nil, fmt.Errorf("shape references unknown material: %d", matId)
		}
		scene.Objects = append(scene.Objects, shape)
	}

	return &scene, nil
}

func LoadSceneJSON(path string) (*rtypes.Scene, error) {
	logrus.Debugf("Loading scene file \"%s\"\n", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	scene, err := ParseScene(data, rregistry.GetCentralRegistry())
	if err != nil {
		return nil, err
	}
	return scene, nil
}
