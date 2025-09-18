package rparser

import (
	"encoding/json"
	"fmt"

	"github.com/piotrwyrw/radia/radia/rimg"
	"github.com/piotrwyrw/radia/radia/rmaterial"
	"github.com/piotrwyrw/radia/radia/robject"
	"github.com/piotrwyrw/radia/radia/rshapes"
	"github.com/piotrwyrw/radia/radia/rtypes"
)

func unmarshalShapeWrapper(data []byte, dst *rtypes.ShapeWrapper) error {
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
		err = sphere.Unmarshal(aux.Object, unmarshalShapeMaterialWrapper)
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

func unmarshalShapeMaterialWrapper(data []byte, dst *rtypes.ShapeMaterialWrapper) error {
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

	var mat rtypes.ShapeMaterial

	switch aux.Name {
	case rtypes.GlassMaterialIdentifier:
		var glass rmaterial.GlassMaterial
		err = json.Unmarshal(aux.Material, &glass)
		mat = &glass
		break
	case rtypes.UniversalMaterialIdentifier:
		var universal rmaterial.UniversalMaterial
		err = json.Unmarshal(aux.Material, &universal)
		mat = &universal
		break
	default:
		return fmt.Errorf("unknown shape material %s", aux.Type)
	}

	if err != nil {
		return err
	}

	*dst = robject.WrapShapeMaterial(mat)

	return nil
}

func unmarshalEnvironmentMaterialWrapper(data []byte, dst *rtypes.EnvironmentMaterialWrapper) error {
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

	var mat rtypes.EnvironmentMaterial

	switch aux.Name {
	case rtypes.SkyMaterialIdentifier:
		var sky rmaterial.Sky
		err = json.Unmarshal(aux.Material, &sky)
		if err != nil {
			break
		}
		if sky.Image == nil {
			return fmt.Errorf("sky material image is missing")
		}
		imgSrc := sky.Image.Source
		sky.Image, err = rimg.RasterFromPNG(imgSrc)
		if err != nil {
			break
		}
		mat = &sky
		break
	default:
		return fmt.Errorf("unknown environment material %s", aux.Type)
	}

	if err != nil {
		return err
	}

	*dst = robject.WrapEnvironmentMaterial(mat)

	return nil
}

func UnmarshalScene(data []byte) (*rtypes.Scene, error) {
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
	err = unmarshalEnvironmentMaterialWrapper(aux.WorldMat, &scene.WorldMat)
	if err != nil {
		return nil, err
	}

	// Parse scene objects
	for _, object := range aux.Objects {
		var shape rtypes.ShapeWrapper
		err = unmarshalShapeWrapper(object, &shape)
		if err != nil {
			return nil, err
		}
		scene.Objects = append(scene.Objects, shape)
	}

	return &scene, nil
}
