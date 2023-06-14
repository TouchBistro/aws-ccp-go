package providers

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/pkg/errors"
)

// DefaultCredsProvider type
type DefaultCredsProvider struct {
	CredsProviderOptions
	awsConfig aws.Config
}

func (p *DefaultCredsProvider) String() string {
	return fmt.Sprintf("Provider Name: %v, Region %v", p.Name, *p.Region)
}

// init initializes the credentials provider using the supplied configuration, or returns an error
func (p *DefaultCredsProvider) init(context context.Context, name string, options ...CredsProviderOptionsFunc) error {

	if len(name) == 0 {
		return ErrInvalidProviderName
	}

	p.Name = name
	p.Region = aws.String(DefaultAWSRegion)
	for _, opt := range options {
		opt(&p.CredsProviderOptions)
	}

	finalOpts := p.LoadOptionFns
	finalOpts = append(finalOpts, config.WithRegion(*p.Region))
	cfg, err := config.LoadDefaultConfig(context, finalOpts...)
	if err != nil {
		return errors.Wrapf(err, "error creating aws config with supplied parameters")
	}

	p.awsConfig = cfg //set the config
	return p.validate()

}

// NewDefaultCredsProvider creates an AWS client provider based on the default credential chain using the supplied options
//
// DefaultCredsProvider is a default Client Provider wrapper; This behaves like the
// underlying AWS SDK client configuration & uses the default credentials chain
// to use environment, shared config or AWS IAM roles in a specified order, determined
// by the AWS SDK itself.
func NewDefaultCredsProvider(context context.Context, name string, options ...CredsProviderOptionsFunc) (*DefaultCredsProvider, error) {

	p := &DefaultCredsProvider{}
	if err := p.init(context, name, options...); err != nil {
		return nil, err
	}

	pmap.put(p.Name, p)
	return p, nil
}

func (p *DefaultCredsProvider) Config() aws.Config { return p.awsConfig }

func (p *DefaultCredsProvider) Key() string { return p.Name }

//internal
func (p *DefaultCredsProvider) validate() error {
	if p.Validation {
		client := sts.NewFromConfig(p.awsConfig)
		_, err := client.GetCallerIdentity(context.Background(), &sts.GetCallerIdentityInput{})

		if err != nil {
			return err
		}

		return nil
	}
	return nil
}
