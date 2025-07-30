package internal

type FirestoreConfig struct {
	Path       string
	ProjectId  string
	DatabaseId string
}

func NewFirestoreConfig(path string, projectId string, databaseId string) *FirestoreConfig {
	return &FirestoreConfig{
		Path:       path,
		ProjectId:  projectId,
		DatabaseId: databaseId,
	}
}
