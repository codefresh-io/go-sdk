package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/codefresh-io/go-sdk/pkg/client"
)

const (
	KubernetesRunnerType = "kubernetes"
)

type (
	// RuntimeEnvironmentAPI declers Codefresh runtime environment API
	RuntimeEnvironmentAPI interface {
		Create(*CreateRuntimeOptions) (*RuntimeEnvironment, error)
		Default(string) (bool, error)
		Delete(string) (bool, error)
		Get(string) (*RuntimeEnvironment, error)
		List() ([]RuntimeEnvironment, error)
		SignCertificate(*SignCertificatesOptions) ([]byte, error)
		Validate(*ValidateRuntimeOptions) error
	}

	runtimeEnvironment struct {
		client *client.CfClient
	}

	RuntimeEnvironment struct {
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
)

// Create - create Runtime-Environment
func (r *runtimeEnvironment) Create(opt *CreateRuntimeOptions) (*RuntimeEnvironment, error) {
	body := map[string]any{
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

	_, err := r.client.RestAPI(context.TODO(), &client.RequestOptions{
		Method: "POST",
		Path:   "/api/custom_clusters/register",
		Body:   body,
	})
	if err != nil {
		return nil, fmt.Errorf("failed creating runtime environment: %w", err)
	}

	re := &RuntimeEnvironment{
		Metadata: RuntimeMetadata{
			Name: fmt.Sprintf("%s/%s", opt.Cluster, opt.Namespace),
		},
	}

	return re, nil
}

func (r *runtimeEnvironment) Default(name string) (bool, error) {
	_, err := r.client.RestAPI(context.TODO(), &client.RequestOptions{
		Method: "PUT",
		Path:   fmt.Sprintf("/api/runtime-environments/default/%s", url.PathEscape(name)),
	})
	if err != nil {
		return false, fmt.Errorf("failed setting default runtime environment: %w", err)
	}

	return true, nil
}

func (r *runtimeEnvironment) Delete(name string) (bool, error) {
	_, err := r.client.RestAPI(context.TODO(), &client.RequestOptions{
		Method: "DELETE",
		Path:   fmt.Sprintf("/api/runtime-environments/%s", url.PathEscape(name)),
	})
	if err != nil {
		return false, fmt.Errorf("failed deleting runtime environment: %w", err)
	}

	return true, nil
}

func (r *runtimeEnvironment) Get(name string) (*RuntimeEnvironment, error) {
	res, err := r.client.RestAPI(context.TODO(), &client.RequestOptions{
		Method: "GET",
		Path:   fmt.Sprintf("/api/runtime-environments/%s", url.PathEscape(name)),
		Query: map[string]any{
			"extend": "false",
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed getting runtime environment: %w", err)
	}

	result := &RuntimeEnvironment{}
	return result, json.Unmarshal(res, result)
}

func (r *runtimeEnvironment) List() ([]RuntimeEnvironment, error) {
	res, err := r.client.RestAPI(context.TODO(), &client.RequestOptions{
		Path:   "/api/runtime-environments",
		Method: "GET",
	})
	if err != nil {
		return nil, fmt.Errorf("failed getting runtime environment list: %w", err)
	}

	result := make([]RuntimeEnvironment, 0)
	return result, json.Unmarshal(res, &result)
}

func (r *runtimeEnvironment) SignCertificate(opt *SignCertificatesOptions) ([]byte, error) {
	res, err := r.client.RestAPI(context.TODO(), &client.RequestOptions{
		Path:   "/api/custom_clusters/signServerCerts",
		Method: "POST",
		Body: map[string]any{
			"reqSubjectAltName": opt.AltName,
			"csr":               opt.CSR,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed signing certificate: %w", err)
	}

	return res, err
}

func (r *runtimeEnvironment) Validate(opt *ValidateRuntimeOptions) error {
	_, err := r.client.RestAPI(context.TODO(), &client.RequestOptions{
		Path:   "/api/custom_clusters/validate",
		Method: "POST",
		Body: map[string]any{
			"clusterName": opt.Cluster,
			"namespace":   opt.Namespace,
		},
	})
	if err != nil {
		return fmt.Errorf("failed validating runtime environment: %w", err)
	}

	return nil
}
