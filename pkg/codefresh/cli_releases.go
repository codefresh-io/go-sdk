package codefresh

import (
	"context"
	"fmt"
)

type (
	ICliReleasesAPI interface {
		GetLatest(ctx context.Context) (string, error)
	}

	CliReleases struct {
		codefresh *codefresh
	}

	graphQlGetLatestReleaseResponse struct {
		Data struct {
			LatestCliRelease string
		}
		Errors []graphqlError
	}
)

func newCliReleasesAPI(codefresh *codefresh) ICliReleasesAPI {
	return &CliReleases{codefresh: codefresh}
}

func (releases *CliReleases) GetLatest(ctx context.Context) (string, error) {
	jsonData := map[string]interface{}{
		"query": `{
			latestCliRelease {
				version
			}
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
