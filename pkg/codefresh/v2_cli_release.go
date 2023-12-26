package codefresh

import (
	"context"
	"fmt"
)

type (
	V2CliReleaseAPI interface {
		GetLatest(ctx context.Context) (string, error)
	}

	v2CliRelease struct {
		codefresh *codefresh
	}

	graphQlGetLatestReleaseResponse struct {
		Data struct {
			LatestCliRelease string
		}
		Errors []graphqlError
	}
)

func (releases *v2CliRelease) GetLatest(ctx context.Context) (string, error) {
	jsonData := map[string]interface{}{
		"query": `{
			latestCliRelease 
		}`,
	}

	res := graphQlGetLatestReleaseResponse{}
	err := releases.codefresh.graphqlAPI(ctx, jsonData, &res)
	if err != nil {
		return "", fmt.Errorf("failed making a graphql API call to runtime: %w", err)
	}

	if len(res.Errors) > 0 {
		return "", graphqlErrorResponse{errors: res.Errors}
	}

	if res.Data.LatestCliRelease == "" {
		return "", fmt.Errorf("failed getting latest release")
	}

	return res.Data.LatestCliRelease, nil
}
