package providers

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/pkg/errors"
)

// EnvironmentCredsProvider type
type EnvironmentCredsProvider struct {
	DefaultCredsProvider
}

type EnvirontmentClientProviderOption func(CredsProvider)

type EnvironmentCredsProviderBuilder struct {
	Provider *EnvironmentCredsProvider
}

func (p *EnvironmentCredsProvider) String() string {
	return fmt.Sprintf("Provider Name: %v, Region: %v,  Access Key ID Envvar:%v, Secret Access Key Envvar: %v, Session Token Envvar: %v, Region Envvar: %v",
		p.Name, *p.Region, *p.AccessKeyIdVar, *p.SecretAccessKeyVar, *p.SessionTokenVar, *p.RegionVar)
}

// init initializes the credentials provider using the supplied configuration, or returns an error
func (p *EnvironmentCredsProvider) init(context context.Context, name string, options ...CredsProviderOptionsFunc) error {

	if len(name) == 0 {
		return ErrInvalidProviderName
	}

	p.Name = name
	p.Region = aws.String(DefaultAWSRegion)
	p.AccessKeyIdVar = aws.String("AWS_ACCESS_KEY_ID")
	p.SecretAccessKeyVar = aws.String("AWS_SECRET_ACCESS_KEY")
	p.SessionTokenVar = aws.String("AWS_SESSION_TOKEN")
	p.RegionVar = aws.String("AWS_REGION")

	for _, opt := range options {
		opt(&p.CredsProviderOptions)
	}

	key := os.Getenv(*p.AccessKeyIdVar)
	if key == "" {
		return ErrInvalidAwsAccessKeyIdEnvValue
	}

	secret := os.Getenv(*p.SecretAccessKeyVar)
	if secret == "" {
		return ErrInvalidSecretAccessKeyEnvValue
	}

	// supplied by WithRegionFrom() or AWS_REGION env value
	region := os.Getenv(*p.RegionVar)
	if region == "" {
		// supplied by WithDefaultRegion() or WithRegion() options
		region = *p.Region
	}
	session := ""
	if p.SessionTokenVar != nil {
		// can be empty, only required for short term AWS principals
		session = os.Getenv(*p.SessionTokenVar)
	}
	staticCredsProvider := credentials.NewStaticCredentialsProvider(key, secret, session)

	finalOpts := p.LoadOptionFns
	finalOpts = append(finalOpts, config.WithRegion(region), config.WithCredentialsProvider(staticCredsProvider))
	cfg, err := config.LoadDefaultConfig(context, finalOpts...)
	if err != nil {
		return errors.Wrapf(err, "error creating aws config with supplied parameters")
	}

	p.awsConfig = cfg //set the config
	return p.validate()

}

// NewEnvironmentCredsProvider creates an AWS client provider based on the static credentials from env vars
// using the supplied options
//
// EnvironmentCredsProvider uses AWS credentials supplied via environment variables. By
// default it uses standard AWS environment variables for static credential. It also allows
// overriding & using other variables name for the three.
//   - AWS_ACCESS_KEY_ID
//   - AWS_SECRET_ACCESS_KEY
//   - AWS_SESSION_TOKEN
//
// AWS Region is set using the following precendence
// 1- If specified with option WithRegionFrom(envvar) & non-empty value is set for that env var
// 2- If a non-empty value exists for env var AWS_REGION
// 3- AWS Regsion set with WithDefaultRegion() or WithRegion(region) options
// 4- The AWsDefaultRegion
//
func NewEnvironmentCredsProvider(context context.Context, name string, options ...CredsProviderOptionsFunc) (*EnvironmentCredsProvider, error) {

	p := &EnvironmentCredsProvider{}
	if err := p.init(context, name, options...); err != nil {
		return nil, err
	}

	pmap.put(p.Name, p)
	return p, nil
}
