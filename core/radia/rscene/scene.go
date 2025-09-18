package rscene

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/piotrwyrw/radia/radia/rparser"
	"github.com/piotrwyrw/radia/radia/rtypes"
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
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	scene, err := rparser.UnmarshalScene(data)
	if err != nil {
		return nil, err
	}
	return scene, nil
}
