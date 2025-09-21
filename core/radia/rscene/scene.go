package rscene

import (
	"encoding/json"
	"math"
	"math/rand/v2"
	"os"
	"os/user"
	"path/filepath"

	"github.com/piotrwyrw/radia/radia/rcolor"
	"github.com/piotrwyrw/radia/radia/rimg"
	"github.com/piotrwyrw/radia/radia/rmaterial"
	"github.com/piotrwyrw/radia/radia/rmath"
	"github.com/piotrwyrw/radia/radia/robject"
	"github.com/piotrwyrw/radia/radia/rshapes"
	"github.com/piotrwyrw/radia/radia/rtypes"
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
		Materials: map[int32]rtypes.ShapeMaterialWrapper{
			0: robject.WrapShapeMaterial(&rmaterial.UniversalMaterial{
				Color:      rcolor.ColorWhite(),
				Emission:   rcolor.ColorBlack(),
				Brightness: 0.0,
				Roughness:  0.5,
			}),
			//1: robject.WrapShapeMaterial(&rmaterial.NormalMaterial{
			//	Brightness: 1.0,
			//}),
			1: robject.WrapShapeMaterial(&rmaterial.UniversalMaterial{
				Color:      rcolor.ColorBlack(),
				Emission:   rcolor.ColorWhite(),
				Brightness: 1,
				Roughness:  0.0,
			}),
			2: robject.WrapShapeMaterial(&rmaterial.MirrorMaterial{
				Color: rcolor.RGB(162, 155, 254),
			}),
			3: robject.WrapShapeMaterial(&rmaterial.MirrorMaterial{
				Color: rcolor.RGB(46, 204, 113),
			}),
		},
		Objects: []rtypes.ShapeWrapper{
			// Ground
			robject.WrapShape(&rshapes.Sphere{
				Center:     rmath.Vec(0, -10000.0, 0),
				Radius:     10000.0,
				MaterialId: 0,
			}),

			// Ceiling
			robject.WrapShape(&rshapes.Sphere{
				Center:     rmath.Vec(0, 10000.0+1, 0),
				Radius:     10000.0,
				MaterialId: 0,
			}),

			// Left Wall
			robject.WrapShape(&rshapes.Sphere{
				Center:     rmath.Vec(10000+0.5, 0.5, 0),
				Radius:     10000,
				MaterialId: 0,
			}),

			// Right Wall
			robject.WrapShape(&rshapes.Sphere{
				Center:     rmath.Vec(-10000-0.5, 0.5, 0),
				Radius:     10000,
				MaterialId: 0,
			}),

			// Back Wall
			robject.WrapShape(&rshapes.Sphere{
				Center:     rmath.Vec(0, 0, 10000+1.5),
				Radius:     10000,
				MaterialId: 0,
			}),

			// Lamp
			robject.WrapShape(&rshapes.Sphere{
				Center:     rmath.Vec(0.0, 1.0, 1.0),
				Radius:     0.2,
				MaterialId: 1,
			}),

			// Mirror
			robject.WrapShape(&rshapes.Sphere{
				Center:     rmath.Vec(0.2, 0.3, 1.0),
				Radius:     0.1,
				MaterialId: 2,
			}),

			// Mirror
			robject.WrapShape(&rshapes.Sphere{
				Center:     rmath.Vec(-.2, 0.3, 1.0),
				Radius:     0.1,
				MaterialId: 3,
			}),
		},
		Camera: rtypes.Camera{
			Location:    rmath.Vec(0, 0.5, 0),
			Facing:      rmath.Vec(0, 0, 1),
			FocalLength: 0.5,
		},
		WorldMat: robject.WrapEnvironmentMaterial(&rmaterial.GradientSky{
			IOR:          1.0,
			Intensity:    0.0,
			ColorSky:     rcolor.RGB(217, 231, 255),
			ColorHorizon: rcolor.RGB(255, 255, 255),
		}),
	}
	return scene
}

func NewRandomTestScene() *rtypes.Scene {
	u, err := user.Current()
	var name string
	if err != nil {
		name = "Unknown"
	} else {
		name = u.Name
	}
	var objects []rtypes.ShapeWrapper

	for i := 0; i < 200; i++ {
		angle := rand.Float64() * math.Pi * 2.0
		distance := rand.Float64() + 0.2
		radius := rand.Float64()*0.02 + 0.01
		x := math.Cos(angle) * distance
		z := math.Sin(angle) * distance
		y := radius
		materialId := int32(0)
		if radius <= 0.012 {
			materialId = 1
		}
		objects = append(objects, robject.WrapShape(&rshapes.Sphere{
			Center:     rmath.Vec(x, y, z),
			Radius:     radius,
			MaterialId: materialId,
		}))
	}

	// Ground Sphere
	gRad := float64(10000)
	objects = append(objects, robject.WrapShape(&rshapes.Sphere{
		Center:     rmath.Vec(0, -gRad, 0),
		Radius:     gRad,
		MaterialId: 0,
	}))

	scene := &rtypes.Scene{
		Metadata: rtypes.SceneMetadata{
			Title:  "New Project",
			Author: name,
		},
		Materials: map[int32]rtypes.ShapeMaterialWrapper{
			0: robject.WrapShapeMaterial(&rmaterial.UniversalMaterial{
				Color:      rcolor.RGB(255, 255, 255),
				Emission:   rcolor.ColorBlack(),
				Brightness: 0.0,
				Roughness:  0.8,
			}),
			1: robject.WrapShapeMaterial(&rmaterial.UniversalMaterial{
				Color:      rcolor.ColorBlack(),
				Emission:   rcolor.RGB(255, 150, 150),
				Brightness: 1.0,
				Roughness:  10,
			}),
		},
		Objects: objects,
		Camera: rtypes.Camera{
			Location:    rmath.Vec(0, 0.02, 0),
			Facing:      rmath.Vec(0, 0, 1),
			FocalLength: 1,
		},
		WorldMat: robject.WrapEnvironmentMaterial(&rmaterial.Sky{
			Image: func() *rimg.Raster {
				raster, _ := rimg.RasterFromPNG("world.png")
				return raster
			}(),
			IOR:       1.0,
			Intensity: 0.2,
			Azimuth:   -2,
		}),
	}

	return scene
}

func SceneMaterialExists(s *rtypes.Scene, id int32) bool {
	_, ok := s.Materials[id]
	return ok
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
