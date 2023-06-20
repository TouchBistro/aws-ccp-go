package providers

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/pkg/errors"
)

// SharedConfigCredsProvider type
type SharedConfigCredsProvider struct {
	DefaultCredsProvider
}

func (p *SharedConfigCredsProvider) String() string {
	return fmt.Sprintf("Provider Name: %v, Region: %v, Credentials File:%v, Config File: %v, Config Profile: %v", p.Name, *p.Region, *p.CredentialsFile, *p.ConfigFile, *p.ConfigProfile)
}

// init initializes the credentials provider using the supplied configuration, or returns an error
func (p *SharedConfigCredsProvider) init(context context.Context, name string, options ...CredsProviderOptionsFunc) error {

	if len(name) == 0 {
		return ErrInvalidProviderName
	}

	homeDir, err := os.UserHomeDir()

	if err != nil {
		return err
	}

	p.Name = name
	p.Region = aws.String(DefaultAWSRegion)
	p.CredentialsFile = aws.String(filepath.Join(homeDir, ".aws", "credentials"))
	p.ConfigFile = aws.String(filepath.Join(homeDir, ".aws", "config"))
	p.ConfigProfile = aws.String("default")

	for _, opt := range options {
		opt(&p.CredsProviderOptions)
	}

	//now initialize aws provider with ctx & shared credentials
	finalOpts := p.LoadOptionFns
	finalOpts = append(finalOpts, config.WithRegion(*p.Region), config.WithSharedCredentialsFiles([]string{*p.CredentialsFile}),
		config.WithSharedConfigFiles([]string{*p.ConfigFile}), config.WithSharedConfigProfile(*p.ConfigProfile))
	cfg, err := config.LoadDefaultConfig(context, finalOpts...)
	if err != nil {
		return errors.Wrapf(err, "error creating aws config with supplied parameters")
	}

	p.awsConfig = cfg //set the config
	return p.validate()

}

// NewSharedConfigCredsProvider creates an AWS client provider based on the AWS SDK shared credentials & config using the supplied options
//
// SharedConfigCredsProvider uses AWS shared configuration files. By default the AWS SDK default shared credentials & config
// files are used. It also allows overriding & loading shard configuration from custom locations
// ~/.aws/credentials
// ~/.aws/cofig
// If not supplied, this provider uses `default` as the default value for the config profile to use.
func NewSharedConfigCredsProvider(context context.Context, name string, options ...CredsProviderOptionsFunc) (*SharedConfigCredsProvider, error) {

	p := &SharedConfigCredsProvider{}
	if err := p.init(context, name, options...); err != nil {
		return nil, err
	}

	pmap.put(p.Name, p)
	return p, nil
}
