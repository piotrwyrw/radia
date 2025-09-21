package rtypes

type Scene struct {
	Metadata  SceneMetadata                  `json:"metadata" ui:"Metadata"`
	Materials map[int32]ShapeMaterialWrapper `json:"materials" ui:"-"`
	Objects   []ShapeWrapper                 `json:"objects" ui:"-"`
	Camera    Camera                         `json:"camera" ui:"Camera"`
	WorldMat  EnvironmentMaterialWrapper     `json:"world" ui:"-"`
}
