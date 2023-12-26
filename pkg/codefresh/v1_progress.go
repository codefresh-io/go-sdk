package codefresh

import "fmt"

type (
	V1ProgressAPI interface {
		Get(string) (*Progress, error)
	}

	v1Progress struct {
		codefresh *codefresh
	}

	Progress struct {
		ID       string   `json:"id"`
		Status   string   `json:"status"`
		Location Location `json:"location"`
	}

	Location struct {
		Type string `json:"type"`
		URL  string `json:"url"`
	}
)

func (p *v1Progress) Get(id string) (*Progress, error) {
	resp, err := p.codefresh.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/progress/%s", id),
		method: "GET",
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	result := &Progress{}
	err = p.codefresh.decodeResponseInto(resp, result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
