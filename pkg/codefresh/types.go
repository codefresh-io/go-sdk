package codefresh

import "net/http"

type (
	// AuthOptions
	AuthOptions struct {

		// Token - Codefresh token
		Token string
	}

	// Options
	ClientOptions struct {
		Auth        AuthOptions
		Debug       bool
		Host        string
		Client      *http.Client
		graphqlPath string
	}

	codefresh struct {
		token       string
		host        string
		graphqlPath string
		client      *http.Client
	}

	requestOptions struct {
		path   string
		method string
		body   interface{}
		qs     interface{}
	}
)
