package codefresh

import (
	"context"
	"fmt"

	"github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
	ICliReleasesAPI interface {
		GetLatest(ctx context.Context) (*model.Release, error)
	}

	CliReleases struct {
		codefresh *codefresh
	}

	graphQlGetReleasesResponse struct {
		Data struct {
			LatestCliRelease model.Release
		}
		Errors []graphqlError
	}
)

func newCliReleaseAPI(codefresh *codefresh) ICliReleasesAPI {
	return &CliReleases{codefresh: codefresh}
}

func (releases *CliReleases) GetLatest(ctx context.Context) (*model.Release, error) {
	jsonData := map[string]interface{}{
		"query": `{
			latestCliRelease {
				version
			}
		}`,
	}

	res := graphQlGetReleasesResponse{}
	err := releases.codefresh.graphqlAPI(ctx, jsonData, &res)
	if err != nil {
		return nil, fmt.Errorf("failed making a graphql API call to runtime: %w", err)
	}

	if len(res.Errors) > 0 {
		return nil, graphqlErrorResponse{errors: res.Errors}
	}

	if &res.Data.LatestCliRelease == (*model.Release)(nil) {
		return nil, fmt.Errorf("failed finding latest release")
	}

	return &res.Data.LatestCliRelease, nil
}
