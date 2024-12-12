package graphql

import (
	"context"
	"fmt"
	platmodel "github.com/codefresh-io/go-sdk/pkg/model/platform"

	"github.com/codefresh-io/go-sdk/pkg/client"
)

type (
	PromotionTemplateAPI interface {
		GetVersionSourceByRuntime(ctx context.Context, app *platmodel.ObjectMeta) (*platmodel.PromotionTemplateShort, error)
	}

	promotionTemplate struct {
		client *client.CfClient
	}
)

func (c *promotionTemplate) GetVersionSourceByRuntime(ctx context.Context, app *platmodel.ObjectMeta) (*platmodel.PromotionTemplateShort, error) {
	query := `
query ($applicationMetadata: Object!) {
    promotionTemplateByRuntime(applicationMetadata: $applicationMetadata) {
    versionSource {
      file
      jsonPath
    }
  }
}`
	variables := map[string]any{
		"applicationMetadata": app,
	}
	versionSource, err := client.GraphqlAPI[platmodel.PromotionTemplateShort](ctx, c.client, query, variables)
	if err != nil {
		return nil, fmt.Errorf("failed to get promotion template: %w", err)
	}

	return &versionSource, nil
}
