package rtypes

type Scene struct {
	Metadata SceneMetadata              `json:"metadata" ui:"Metadata"`
	Objects  []ShapeWrapper             `json:"objects" ui:"-"`
	Camera   Camera                     `json:"camera" ui:"Camera"`
	WorldMat EnvironmentMaterialWrapper `json:"world" ui:"-"`
}
