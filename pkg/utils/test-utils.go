package utils

import (
	"os"
	"path/filepath"

	"github.com/codefresh-io/go-sdk/pkg/client"
)

func NewClientFromCurrentContext() *client.CfClient {
	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(homeDir, ".cfconfig")
	authContext, _ := ReadAuthContext(path, "")
	return client.NewCfClient(authContext.URL, authContext.Token, "", nil)
}
