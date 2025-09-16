package rscene

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Scene struct {
	Objects  []Shape
	Camera   Camera
	WorldMat EnvironmentMaterial
}

func (s *Scene) SaveJSON(path string) error {
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
