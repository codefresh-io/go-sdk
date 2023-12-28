package utils

import (
	"fmt"
	"os"

	"sigs.k8s.io/yaml"
)

type (
	CFContext struct {
		Type        string `json:"type"`
		Name        string `json:"name"`
		URL         string `json:"url"`
		Token       string `json:"token"`
		Beta        bool   `json:"beta"`
		OnPrem      bool   `json:"onPrem"`
		ACLType     string `json:"acl-type"`
		UserID      string `json:"user-id"`
		AccountID   string `json:"account-id"`
		Expires     int    `json:"expires"`
		UserName    string `json:"user-name"`
		AccountName string `json:"account-name"`
	}

	CFConfig struct {
		Contexts       map[string]*CFContext `json:"contexts"`
		CurrentContext string                `json:"current-context"`
	}
)

func ReadAuthContext(path string, name string) (*CFContext, error) {
	config, err := getCFConfig(path)
	if err != nil {
		return nil, err
	}

	var context *CFContext
	if name != "" {
		context = config.Contexts[name]
	} else {
		context = config.Contexts[config.CurrentContext]
	}

	return context, nil
}

func getCFConfig(path string) (*CFConfig, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading file\n")
		return nil, err
	}

	config := CFConfig{}
	err = yaml.Unmarshal(content, &config)
	if err != nil {
		fmt.Printf("Error unmarshaling content\n")
		fmt.Println(err.Error())
		return nil, err
	}

	return &config, nil
}
