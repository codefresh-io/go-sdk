package codefresh

import (
  "context"
  "fmt"

  platmodel "github.com/codefresh-io/go-sdk/pkg/codefresh/model"
)

type (
  IAccountV2API interface {
    UpdateCsdpSettings(ctx context.Context, gitProvider platmodel.GitProviders, gitApiUrl, sharedConfigRepo string) error
  }

  accountV2 struct {
    codefresh *codefresh
  }
)

func newAccountV2API(codefresh *codefresh) IAccountV2API {
  return &accountV2{codefresh: codefresh}
}

func (c *accountV2) UpdateCsdpSettings(ctx context.Context, gitProvider platmodel.GitProviders, gitApiUrl, sharedConfigRepo string) error {
  jsonData := map[string]interface{}{
    "query": `
      mutation updateCsdpSettings($gitProvider: GitProviders!, $gitApiUrl: String!, $sharedConfigRepo: String!) {
        updateCsdpSettings(gitProvider: $gitProvider, gitApiUrl: $gitApiUrl, sharedConfigRepo: $sharedConfigRepo)
      }
    `,
    "variables": map[string]interface{}{
      "gitProvider":      gitProvider,
      "gitApiUrl":        gitApiUrl,
      "sharedConfigRepo": sharedConfigRepo,
    },
  }
  res := &graphqlVoidResponse{}
  err := c.codefresh.graphqlAPI(ctx, jsonData, res)
  if err != nil {
    return fmt.Errorf("failed making a graphql API call to update csdp settings: %w", err)
  }

  if len(res.Errors) > 0 {
    return graphqlErrorResponse{errors: res.Errors}
  }

  return nil
}
