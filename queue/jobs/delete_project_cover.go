package jobs

type DeleteProjectCoverArgs struct {
	ProjectID int32  `json:"project_id"`
	Key       string `json:"key"`
}

func (DeleteProjectCoverArgs) Kind() string { return "delete_project_cover" }
