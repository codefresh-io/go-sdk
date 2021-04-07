package codefresh

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var clientOptions = ClientOptions{
	Auth: AuthOptions{
		Token: "",
	},
	Debug: true,
	Host:  "https://g.codefresh.io",
}

func TestProject(t *testing.T) {

	cf := New(&clientOptions)
	projects, err := cf.Projects().List()
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, len(projects), 0)
}
