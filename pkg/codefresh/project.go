package codefresh

type (
	IProjectAPI interface {
		List() ([]*Project, error)
	}
	project struct {
		codefresh *codefresh
	}
	Project struct {
		ProjectName    string `json:"projectName"`
		PipelineNumber int    `json:"pipelineNumber"`
	}
	getProjectResponse struct {
		Total    int        `json:"limit"`
		Projects []*Project `json:"projects"`
	}
)

func newProjectAPI(codefresh *codefresh) IProjectAPI {
	return &project{codefresh}
}

func (p *project) List() ([]*Project, error) {
	r := &getProjectResponse{}

	resp, err := p.codefresh.requestAPI(&requestOptions{
		path:   "/api/projects",
		method: "GET",
	})
	if err != nil {
		return nil, err
	}
	err = p.codefresh.decodeResponseInto(resp, r)
	if err != nil {
		return nil, err
	}
	return r.Projects, nil
}
