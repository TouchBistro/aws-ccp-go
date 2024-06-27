package awscredsprovider

import (
	"context"

	"github.com/TouchBistro/aws-ccp-go/providers"
)

// Creates a provider named "default" of the type `DefaultCredsProvider`, which uses the AWS SDK default credentials chain.
func init() {
	_, _ = providers.NewDefaultCredsProvider(context.Background(), providers.DefaultCredsProviderName, providers.WithDefaultRegion())
}
