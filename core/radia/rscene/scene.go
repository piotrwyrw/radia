package rscene

import (
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"

	"github.com/piotrwyrw/radia/radia/rcolor"
	"github.com/piotrwyrw/radia/radia/rmaterial"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/robject"
	"github.com/piotrwyrw/radia/radia/rparser"
	"github.com/piotrwyrw/radia/radia/rregistry"
	"github.com/piotrwyrw/radia/radia/rshapes"
	"github.com/piotrwyrw/radia/radia/rtypes"
	"github.com/sirupsen/logrus"
)

func NewBlankScene() *rtypes.Scene {
	u, err := user.Current()
	var name string
	if err != nil {
		name = "Unknown"
	} else {
		name = u.Name
	}
	scene := &rtypes.Scene{
		Metadata: rtypes.SceneMetadata{
			Title:  "New Project",
			Author: name,
		},
		Objects: []rtypes.ShapeWrapper{
			robject.WrapShape(&rshapes.Sphere{
				Center: rmath.Vec(0, 0, 2),
				Radius: 0.5,
				Material: robject.WrapShapeMaterial(&rmaterial.UniversalMaterial{
					Color:      rcolor.RGB(244, 162, 137),
					Emission:   rcolor.ColorBlack(),
					Brightness: 0.0,
					Roughness:  0,
				}),
			}),
		},
		Camera: rtypes.Camera{
			Location:    rmath.Vec(0, 0, 0),
			Facing:      rmath.Vec(0, 0, 1),
			FocalLength: 0.5,
		},
		WorldMat: robject.WrapEnvironmentMaterial(&rmaterial.GradientSky{
			IOR:          1.0,
			Intensity:    2.0,
			ColorHorizon: rcolor.RGB(255, 190, 118),
			ColorSky:     rcolor.RGB(19, 15, 64),
		}),
	}
	return scene
}

func SaveSceneJSON(s *rtypes.Scene, path string) error {
	dir := filepath.Dir(path)
	e := os.MkdirAll(dir, 0777)
	if e != nil {
		return e
	}

	f, e := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0777)
	if e != nil {
		return e
	}
	defer f.Close()

	e = json.NewEncoder(f).Encode(s)
	if e != nil {
		return e
	}

	return nil
}

func LoadSceneJSON(path string) (*rtypes.Scene, error) {
	logrus.Debugf("Loading scene file \"%s\"\n", path)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	scene, err := rparser.ParseScene(data, rregistry.GetCentralRegistry())
	if err != nil {
		return nil, err
	}
	return scene, nil
}
