package cmd

import (
	"context"

	"github.com/TouchBistro/aws-ccp-go/providers"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// default value for all string flags
const defaultStrValue string = ""

// default value for all bool flgas
const defaultBoolValue bool = false

const (
	fl_region              string = "region"
	fl_profile             string = "profile"
	fl_credsFile           string = "creds-file"
	fl_configFile          string = "config-file"
	fl_env                 string = "env"
	fl_accessKeyIdFrom     string = "access-key-id-from"
	fl_secretAccessKeyFrom string = "secret-access-key-from"
	fl_sessionTokenFrom    string = "session-token-from"
	fl_roleArn             string = "role-arn"
)

// When aws-ccp-go clients import this piackage with a blank import, the following command line flags are parsed & used to
// configure the "default" creds provider.
//
// Importing this package is non-intrusive, as it should not affect any command-line processing for the client code.
// however, some third-party packages will throw an error if the flags supplied at runtime to configure aws-ccp-go are not defined
// & registered with those packages. So any subset of flags that the client code intends on using may also require to be configured
// with the client code's command-line processing library.
//
// # General
//
// region :  When supplied, configures the default AWS region to use. 'us-east-1' is the hard-coded default currently.
//
//	(This flag only impacts the auto-initialized `default` provider)
//
// profile : When supplied, the auto-initialized default provider will be configured using a `SharedCredsProvider`
// and will use the supplied value as the profile name. This value overrides any other option on the command-line.
//
// # Shared Credential Configuration Tuning
//
// if `profile` flag is supplied, the following two options can be used to configure the AWS config & credential file
// location
//
//	creds-file :  When supplied, the default provider uses the supplied file as the AWS Credentials file for shared config
//			      else default `~/.aws/credentials` is used.
//
//	config-file :  When supplied,the default provider uses the supplied file as the AWS Config file for shared config
//		           else default `~/.aws/config` is used.
//
//	env : When supplied, explicitly indicates the default provider should use the EnvironmentCredsProvider & the default
//		  env var names to get the access_key_id, secret_access_key & session_token values. if supplied with `profile` flag
//		  the profile flag takes precendece & the auto-initialized default provider will use the SharedConfigCredsProvider
//		  instead
//
// # Environment Variables Configuration Tuning
//
//		 access-key-id-from: When `env` flag is supplied, the valud for this flag specifies the env var name to use for fetching
//			                    the AWS Acces Key ID. If not supplied, AWS_ACCESS_KEY_ID is the deafult value used
//
//		 secret-access-key-from: When `env` flag is supplied, the valud for this flag specifies the env var name to use for fetching
//			                       the AWS Secret Acces Key. If not supplied, AWS_SECRET_ACCESS_KEY_ID is the deafult value used
//
//		 session-token-from: When `env` flag is supplied, the valud for this flag specifies the env var name to use for fetching
//			                    the AWS Session Token. If not supplied, AWS_SESSION_TOKEN is the deafult value used.
//
//
//	 # Assume Role
//
//	  role-arn: When this value is supplied, the `default`` provider uses the `AssumeRoleCredsProvider` & the supplied role arn as the role
//	            to assume. The base provider to use for
//	            This option works in conjunction will all others above, using the same rules, however, after the `default`
//	            provider has been initialized, it is then used to assume the role arn supplied by this flag.
func init() {

	// region
	addHiddenStringFlag(fl_region, defaultStrValue, "aws default region for the default provider")

	// profile
	addHiddenStringFlag(fl_profile, defaultStrValue, "aws profile for shared credentials")
	addHiddenStringFlag(fl_credsFile, defaultStrValue, "the credentials file path")
	addHiddenStringFlag(fl_configFile, defaultStrValue, "the config file path")

	// env
	addHiddenBoolFlag(fl_env, defaultBoolValue, "use env")
	addHiddenStringFlag(fl_accessKeyIdFrom, defaultStrValue, "env var to read access key id from")
	addHiddenStringFlag(fl_secretAccessKeyFrom, defaultStrValue, "env var to read secret access key from")
	addHiddenStringFlag(fl_sessionTokenFrom, defaultStrValue, "env var to read session toekn from")

	// assume role
	addHiddenStringFlag(fl_roleArn, defaultStrValue, "aws rold arn to assume")

	// ignore unknown flags
	pflag.CommandLine.ParseErrorsWhitelist.UnknownFlags = true
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	region := viper.GetString(fl_region)
	profile := viper.GetString(fl_profile)
	configFile := viper.GetString(fl_configFile)
	credsFile := viper.GetString(fl_credsFile)
	isEnv := viper.GetBool(fl_env)
	accessKeyId := viper.GetString(fl_accessKeyIdFrom)
	secretAccessKey := viper.GetString(fl_secretAccessKeyFrom)
	sessionToken := viper.GetString(fl_sessionTokenFrom)
	roleArn := viper.GetString(fl_roleArn)

	defaultRegion := providers.DefaultAWSRegion
	if supplied(region) {
		defaultRegion = region
	}

	fn := make([]providers.CredsProviderOptionsFunc, 0)
	fn = append(fn, providers.WithRegion(defaultRegion))

	if supplied(profile) {
		fn = append(fn, providers.WithConfigProfile(profile))
		if supplied(credsFile) {
			fn = append(fn, providers.WithCredentialsFile(credsFile))
		}
		if supplied(configFile) {
			fn = append(fn, providers.WithConfigFile(configFile))
		}
		_, _ = providers.NewSharedConfigCredsProvider(context.Background(), providers.DefaultCredsProviderName, fn...)
	} else if isEnv {
		if supplied(accessKeyId) {
			fn = append(fn, providers.WithAccessKeyIdFrom(accessKeyId))
		}
		if supplied(secretAccessKey) {
			fn = append(fn, providers.WithSecretAccessKeyFrom(secretAccessKey))
		}
		if supplied(sessionToken) {
			fn = append(fn, providers.WithSessionTokenFrom(sessionToken))
		}
		_, _ = providers.NewEnvironmentCredsProvider(context.Background(), providers.DefaultCredsProviderName, fn...)
	} else {
		// use the default creds provider
		_, _ = providers.NewDefaultCredsProvider(context.Background(), providers.DefaultCredsProviderName, fn...)
	}

	// now check if a role arn is supplied
	if supplied(roleArn) {
		fn = make([]providers.CredsProviderOptionsFunc, 0)
		fn = append(fn, providers.WithRegion(defaultRegion))
		fn = append(fn, providers.WithBaseCredsProviderName(providers.DefaultCredsProviderName))
		fn = append(fn, providers.WithRoleArn(roleArn))
		_, _ = providers.NewAssumeRoleCredsProvider(context.Background(), providers.DefaultCredsProviderName, fn...)
	}
}

// addHiddenStringFlag defines a pflag string variable that is hidden from help messages
func addHiddenStringFlag(name, def, usage string) {
	pflag.String(name, def, usage)
	pflag.CommandLine.MarkHidden(name)
}

// addHiddenBoolFlag defines a pflag bool variable that is hidden from help messages
func addHiddenBoolFlag(name string, def bool, usage string) {
	pflag.Bool(name, def, usage)
	pflag.CommandLine.MarkHidden(name)
}

func supplied(val string) bool {
	return val != "" && len(val) > 0
}
