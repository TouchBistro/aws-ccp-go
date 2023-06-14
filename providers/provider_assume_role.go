package providers

import (
	"context"
	"fmt"

	"github.com/TouchBistro/aws-ccp-go/util"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type AssumeRoleCredsProvider struct {
	DefaultCredsProvider
}

func (p *AssumeRoleCredsProvider) String() string {
	return fmt.Sprintf("Provider Name: %v, Region %v, Role ARN: %#v, AWS Account ID: %#v, AWS Role Name: %#v",
		p.Name, util.Coalesce(p.Region, "nil"), util.Coalesce(p.RoleArn, "nil"), util.Coalesce(p.AccountId, "nil"), util.Coalesce(p.RoleName, "nil"))
}

// init initializes the credentials provider using the supplied configuration, or returns an error
func (p *AssumeRoleCredsProvider) init(context context.Context, name string, options ...CredsProviderOptionsFunc) error {

	if len(name) == 0 {
		return ErrInvalidProviderName
	}

	p.Name = name
	p.Region = aws.String(DefaultAWSRegion)

	for _, opt := range options {
		opt(&p.CredsProviderOptions)
	}

	//set base config
	var cfg aws.Config
	if p.BaseCredsProvider != nil {
		cfg = p.BaseCredsProvider.Config()
	} else if p.BaseCredsProviderName != nil {
		if pr, ok := pmap.get(*p.BaseCredsProviderName); ok {
			cfg = pr.Config()
		} else {
			return ErrInvalidBaseProviderConfig
		}
	} else {
		var err error
		finalOpts := p.LoadOptionFns //awsConfigOptFns
		finalOpts = append(finalOpts, config.WithRegion(*p.Region))
		cfg, err = config.LoadDefaultConfig(context, finalOpts...)
		if err != nil {
			return err
		}
	}

	// get the role arn to assume for credentials, if roleArn is supplied, use that,
	// else if both roleName & accountId are supplied, use that;
	var role *string
	if p.RoleArn != nil {
		role = p.RoleArn
	} else if p.RoleName != nil && p.AccountId != nil {
		role = aws.String(fmt.Sprintf("arn:aws:iam::%v:role/%v", *p.AccountId, *p.RoleName))
	}

	// if role is provided, either as arn or role name/account id tuple, assume that role & use sts credentials
	// if no role is supplied, uust use the cfg from the base provider as-is
	var stsSvc *sts.Client
	var creds *stscreds.AssumeRoleProvider
	if role != nil {
		stsSvc = sts.NewFromConfig(cfg, func(o *sts.Options) { o.Region = *p.Region })
		creds = stscreds.NewAssumeRoleProvider(stsSvc, *role, func(aro *stscreds.AssumeRoleOptions) {
		})
		cfg.Credentials = aws.NewCredentialsCache(creds)
	} else {
		p.Region = &cfg.Region //use base provider region
	}

	p.awsConfig = cfg
	p.awsConfig.Region = *p.Region
	return p.validate()
}

// NewAssumeRoleCredsProvider creates an AWS client provider with base credentials and an
// assumed role using the supplied options
//
// AssumeRoleCredsProvider is an extension of the supplied CredsProvider, but additionally
// uses a role arn (or an aws accountid & role name) to assume that role;
//
// When using this provider, the credentials obtained using base `CredsProvider` are used to
// assume the supplied role & the stscreds are used to obtain the AWS credentials for this provider.
//
// If no base CredsProvider is supplied, a DefaultCredsProvider is initialized using the default
// credentials chain & supplied AWS region.
//
// If no RoleArn (or RoleName and AccountID) is supplied, the base CredsProvider credentials are
// used for this provider's aws.Config.
//
// AWS Region supplied as options for this provider are used as the final AWS region in the
// aws.Config.
//
// The default or (supplied base) credentials must have the required permissions to be able
// to assume the role
func NewAssumeRoleCredsProvider(context context.Context, name string, options ...CredsProviderOptionsFunc) (*AssumeRoleCredsProvider, error) {

	p := &AssumeRoleCredsProvider{}
	if err := p.init(context, name, options...); err != nil {
		return nil, err
	}
	pmap.put(p.Name, p)
	return p, nil
}
