package codefresh

import (
	"fmt"
	"io"
	"net/url"
	"time"
)

const (
	KubernetesRunnerType = "kubernetes"
)

type (
	// V1RuntimeEnvironmentAPI declers Codefresh runtime environment API
	V1RuntimeEnvironmentAPI interface {
		Create(*CreateRuntimeOptions) (*v1RuntimeEnvironment, error)
		Default(string) (bool, error)
		Delete(string) (bool, error)
		Get(string) (*v1RuntimeEnvironment, error)
		List() ([]v1RuntimeEnvironment, error)
		SignCertificate(*SignCertificatesOptions) ([]byte, error)
		Validate(*ValidateRuntimeOptions) error
	}

	v1RuntimeEnvironment struct {
		Version               int                   `json:"version"`
		Metadata              RuntimeMetadata       `json:"metadata"`
		Extends               []string              `json:"extends"`
		Description           string                `json:"description"`
		AccountID             string                `json:"accountId"`
		RuntimeScheduler      RuntimeScheduler      `json:"runtimeScheduler"`
		DockerDaemonScheduler DockerDaemonScheduler `json:"dockerDaemonScheduler"`
		Status                struct {
			Message   string    `json:"message"`
			UpdatedAt time.Time `json:"updated_at"`
		} `json:"status"`
	}

	RuntimeScheduler struct {
		Cluster struct {
			ClusterProvider struct {
				AccountID string `json:"accountId"`
				Selector  string `json:"selector"`
			} `json:"clusterProvider"`
			Namespace string `json:"namespace"`
		} `json:"cluster"`
		UserAccess bool `json:"userAccess"`
		Pvcs       struct {
			Dind struct {
				StorageClassName string `yaml:"storageClassName"`
			} `yaml:"dind"`
		} `yaml:"pvcs"`
	}

	DockerDaemonScheduler struct {
		Cluster struct {
			ClusterProvider struct {
				AccountID string `json:"accountId"`
				Selector  string `json:"selector"`
			} `json:"clusterProvider"`
			Namespace string `json:"namespace"`
		} `json:"cluster"`
		UserAccess bool `json:"userAccess"`
	}

	RuntimeMetadata struct {
		Agent        bool   `json:"agent"`
		Name         string `json:"name"`
		ChangedBy    string `json:"changedBy"`
		CreationTime string `json:"creationTime"`
	}

	CreateRuntimeOptions struct {
		Cluster            string
		Namespace          string
		HasAgent           bool
		StorageClass       string
		RunnerType         string
		DockerDaemonParams string
		NodeSelector       map[string]string
		Annotations        map[string]string
	}

	ValidateRuntimeOptions struct {
		Cluster   string
		Namespace string
	}

	SignCertificatesOptions struct {
		AltName string
		CSR     string
	}

	CreateResponse struct {
		Name string
	}

	runtimeEnvironment struct {
		codefresh *codefresh
	}
)

// Create - create Runtime-Environment
func (r *runtimeEnvironment) Create(opt *CreateRuntimeOptions) (*v1RuntimeEnvironment, error) {
	re := &v1RuntimeEnvironment{
		Metadata: RuntimeMetadata{
			Name: fmt.Sprintf("%s/%s", opt.Cluster, opt.Namespace),
		},
	}
	body := map[string]interface{}{
		"clusterName":        opt.Cluster,
		"namespace":          opt.Namespace,
		"storageClassName":   opt.StorageClass,
		"runnerType":         opt.RunnerType,
		"dockerDaemonParams": opt.DockerDaemonParams,
		"nodeSelector":       opt.NodeSelector,
		"annotations":        opt.Annotations,
	}
	if opt.HasAgent {
		body["agent"] = true
	}

	resp, err := r.codefresh.requestAPI(&requestOptions{
		path:   "/api/custom_clusters/register",
		method: "POST",
		body:   body,
	})

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode < 400 {
		return re, nil
	}

	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return nil, fmt.Errorf("Error during runtime environment creation, error: %s", string(buffer))
}

func (r *runtimeEnvironment) Default(name string) (bool, error) {
	path := fmt.Sprintf("/api/runtime-environments/default/%s", url.PathEscape(name))
	resp, err := r.codefresh.requestAPI(&requestOptions{
		path:   path,
		method: "PUT",
	})
	if err != nil {
		return false, err
	}

	defer resp.Body.Close()
	if resp.StatusCode == 201 {
		return true, nil
	} else {
		res, err := r.codefresh.getBodyAsString(resp)
		if err != nil {
			return false, err
		}

		return false, fmt.Errorf("Unknown error: %v", res)
	}
}

func (r *runtimeEnvironment) Delete(name string) (bool, error) {
	resp, err := r.codefresh.requestAPI(&requestOptions{
		path:   fmt.Sprintf("/api/runtime-environments/%s", url.PathEscape(name)),
		method: "DELETE",
	})
	if err != nil {
		return false, err
	}

	if resp.StatusCode < 400 {
		return true, nil
	}
	body, err := r.codefresh.getBodyAsString(resp)
	if err != nil {
		return false, err
	}
	return false, fmt.Errorf(body)
}

func (r *runtimeEnvironment) Get(name string) (*v1RuntimeEnvironment, error) {
	path := fmt.Sprintf("/api/runtime-environments/%s", url.PathEscape(name))
	resp, err := r.codefresh.requestAPI(&requestOptions{
		path:   path,
		method: "GET",
		qs: map[string]string{
			"extend": "false",
		},
	})

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	re := &v1RuntimeEnvironment{}
	r.codefresh.decodeResponseInto(resp, re)
	return re, nil
}

func (r *runtimeEnvironment) List() ([]v1RuntimeEnvironment, error) {
	resp, err := r.codefresh.requestAPI(&requestOptions{
		path:   "/api/runtime-environments",
		method: "GET",
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	emptySlice := make([]v1RuntimeEnvironment, 0)
	r.codefresh.decodeResponseInto(resp, emptySlice)
	if err != nil {
		return nil, err
	}

	return emptySlice, err
}

func (r *runtimeEnvironment) SignCertificate(opt *SignCertificatesOptions) ([]byte, error) {
	body := map[string]interface{}{
		"reqSubjectAltName": opt.AltName,
		"csr":               opt.CSR,
	}
	resp, err := r.codefresh.requestAPI(&requestOptions{
		path:   "/api/custom_clusters/signServerCerts",
		method: "POST",
		body:   body,
	})
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	return r.codefresh.getBodyAsBytes(resp)
}

func (r *runtimeEnvironment) Validate(opt *ValidateRuntimeOptions) error {
	body := map[string]interface{}{
		"clusterName": opt.Cluster,
		"namespace":   opt.Namespace,
	}
	resp, err := r.codefresh.requestAPI(&requestOptions{
		path:   "/api/custom_clusters/validate",
		method: "POST",
		body:   body,
	})
	defer resp.Body.Close()
	return err
}
