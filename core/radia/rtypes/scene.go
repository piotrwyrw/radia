package rtypes

type Scene struct {
	Metadata SceneMetadata              `json:"metadata"`
	Objects  []ShapeWrapper             `json:"objects"`
	Camera   Camera                     `json:"camera"`
	WorldMat EnvironmentMaterialWrapper `json:"world"`
}
