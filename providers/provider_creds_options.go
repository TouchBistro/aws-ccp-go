package providers

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// CredsProviderOptionsFunc is a type alias for CredsProviderOptions functional option
type CredsProviderOptionsFunc func(*CredsProviderOptions)

// CredsProviderOptions are a set of options that are valid for `CredProviders` types
// Not all options are used for all `CredProviders`.
type CredsProviderOptions struct {

	// The following options apply to all CredsProviders

	// Name for this provider. This name would be used to store & retrieve the provider
	// using the providers.Get() function.
	Name string

	// The AWS Region to set in the config.
	Region *string

	// Is aws.Config to be validated after initialization. Default is false.
	Validation bool

	// Additional LoadOptions to pass to config.LoadDefaultConfig(...) AWS SDK
	// API.
	LoadOptionFns []func(*config.LoadOptions) error

	// The following options only apply to the EnvironmentCredsProvider

	// The Environment Variable name to fetch the AWS Access Key Id. When not
	// supplied, AWS_ACCESS_KEY_ID is the default.
	AccessKeyIdVar *string

	// The Environment Variable name to fetch the AWS Secret Access Key. When
	// not supplied, AWS_SECRET_ACCESS_KEY is the default.
	SecretAccessKeyVar *string

	// The Environment Variable name to fetch the AWS Session Token. When
	//not supplied AWS_SESSION_TOKEN is the default.
	SessionTokenVar *string

	// The Environment Variable name to fetch the AWS Region from. When
	//no supplied AWS_REGION is the default.
	RegionVar *string

	// The following options only apply to the SharedConfigCredsProvider

	// The AWS credentials file to use. If not supplied, the default is
	// ~/.aws/credentials
	CredentialsFile *string

	// The AWS config file to use. If not supplied, the default is
	// ~/.aws/config
	ConfigFile *string

	// The config profile to use. If not supplied, the default is
	// the profile named `default`
	ConfigProfile *string

	// The following options only apply to the AssumeRoleCredsProvider

	// The Role ARN to assume for the final credentials to use. When Role
	// ARN is supplied, the AccountID & RoleName are ignored.
	RoleArn *string

	// The Account ID for the role to assume for the final credentials to
	// use. RoleName must also be supplied. If RoleArn is supplied, it
	// takes precedence
	AccountId *string

	// The Role Name to assume for the final credentials. AccountId must
	// also be supplied. If RoleArn is supplied, it takes precedence
	RoleName *string

	// The base CredsProvider name, to use & build the STS client for assuming
	// the role for the final credentials. If the named provider does
	// not exists this value is ignore. If BaseCredsProvider value is
	// supplied it takes precedence.
	BaseCredsProviderName *string

	// The base CredsProvider to use & build the STS client for assuming
	// the role for the final credentials. If the BaseCredsProviderName is
	// also supplied, it is ignored & this CredsProvider takes precedence
	BaseCredsProvider CredsProvider
}

// WithDefaultRegion sets `providers.AWSDefaultRegion` (us-east-1) as the AWS Region to use
// by the ßnderlying aws.onfig
func WithDefaultRegion() CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.Region = aws.String(DefaultAWSRegion)
	}
}

// WithRegion sets the supplied region as the AWS Region to use by the underlying
// aws.Config
func WithRegion(region string) CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.Region = aws.String(region)
	}
}

// WithConfigLoadOptFns supplies functional options to pass additional configuration options
// directly to underlying calls to `config.LoadDefaultConfig()`
func WithConfigLoadOptFns(optFns ...func(*config.LoadOptions) error) CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.LoadOptionFns = optFns
	}
}

// ValidateProvider turns on credential validation. This acts as an early failure check.
// The NewXXXCredsProvider() builder functions fails with a no-nil error if the credentials
// are invalid.
//
// If invalid credentials are not validated at this stage, any API operations using an AWS
// SDK client generated with this provider will result in errors due to failure to sign requests
// properly.
//
// The validation step performs an `sts:GetCallerIdentity()` operation which does not require
// any specific permissions.
func ValidateProvider() CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.Validation = true
	}
}

// WithAccessKeyIdFrom specify the environemt variable to use to read access key id
func WithAccessKeyIdFrom(envVarKey string) CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.AccessKeyIdVar = aws.String(envVarKey)
	}
}

// WithSecretAccessKeyFrom specify the environment variable to use to read secret access key
func WithSecretAccessKeyFrom(envVarKey string) CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.SecretAccessKeyVar = aws.String(envVarKey)
	}
}

// WithSessionTokenFrom specify the environment variable to use to read session token
func WithSessionTokenFrom(envVarKey string) CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.SessionTokenVar = aws.String(envVarKey)
	}
}

// WithRegionFrom specify the environment variable to use to read aws region
func WithRegionFrom(envVarKey string) CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.RegionVar = aws.String(envVarKey)
	}
}

// WithCredentialsFile specify path for the credentials file to use
func WithCredentialsFile(path string) CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.CredentialsFile = aws.String(path)
	}
}

// WithConfigFile specify config path for the config file to use
func WithConfigFile(path string) CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.ConfigFile = aws.String(path)
	}
}

// WithConfigProfile specify config path for the config file to use
func WithConfigProfile(profile string) CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.ConfigProfile = aws.String(profile)
	}
}

// WithRoleArn specify the role arn to assume; if supplied account id and role name are ignored
func WithRoleArn(arn string) CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.RoleArn = aws.String(arn)
	}
}

// WithAccountId specify the aws account Id for the role to assume; must also specify role name
func WithAccountId(accountid string) CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.AccountId = aws.String(accountid)
	}
}

// WithRoleName specify the role name for the role to assume; must also specify account id
func WithRoleName(name string) CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.RoleName = aws.String(name)
	}
}

// WithBaseCredsProvideName specify the name of the existing creds provider to use as the baseline provider
// to assume the role supplied. These credentials must be for a princpal that has sts:assumeRole
// permissions on the supplied role arn; If a base CredsProvider is also supplied, that option takes
// precendence over this.
func WithBaseCredsProviderName(name string) CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.BaseCredsProviderName = aws.String(name)
	}
}

// WithBaseCredsProvider supply a creds provider to use as the baseline provider to assume the role supplied.
// These credentials must be for a princpal that has sts:assumeRole permissions on the supplied role arn;
// This option takes precendence over the creds provider name.
func WithBaseCredsProvider(base CredsProvider) CredsProviderOptionsFunc {
	return func(provider *CredsProviderOptions) {
		provider.BaseCredsProvider = base
	}
}
