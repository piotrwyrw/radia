package rscene

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/piotrwyrw/radia/radia/rparser"
	"github.com/piotrwyrw/radia/radia/rregistry"
	"github.com/piotrwyrw/radia/radia/rtypes"
	"github.com/sirupsen/logrus"
)

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
