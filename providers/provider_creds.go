package providers

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
)

// CredsProvider is the interface for all CredsProvider types
type CredsProvider interface {
	init(context.Context, string, ...CredsProviderOptionsFunc) error
	validate() error

	Key() string
	Config() aws.Config
}
