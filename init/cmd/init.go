package cmd

import (
	"context"
	"flag"

	"github.com/TouchBistro/aws-ccp-go/providers"
)

// When aws-ccp-go clients import this piackage with a blank import, the following command line flags are parsed & used to
// configure a "default" creds provider.
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
//	 access-key-id-from: When `env` flag is supplied, the valud for this flag specifies the env var name to use for fetching
//		                    the AWS Acces Key ID. If not supplied, AWS_ACCESS_KEY_ID is the deafult value used
//
//	 secret-access-key-from: When `env` flag is supplied, the valud for this flag specifies the env var name to use for fetching
//		                       the AWS Secret Acces Key. If not supplied, AWS_SECRET_ACCESS_KEY_ID is the deafult value used
//
//	 session-token-from: When `env` flag is supplied, the valud for this flag specifies the env var name to use for fetching
//		                    the AWS Session Token. If not supplied, AWS_SESSION_TOKEN is the deafult value used.
func init() {

	// region
	region := flag.String("region", "", "supply the aws default region")

	// profile
	profile := flag.String("profile", "", "supply the profile")
	credsFile := flag.String("creds-file", "", "env var to read the credentials file path")
	configFile := flag.String("config-file", "", "env var to read the config file path")

	// env
	isEnv := flag.Bool("env", false, "use env")
	accessKeyId := flag.String("access-key-id-from", "", "env var to read access key id from")
	secretAccessKey := flag.String("secret-access-key-from", "", "env var to read secret access key from")
	sessionToken := flag.String("session-name-from", "", "env var to read session toekn from")

	flag.Parse()

	defaultRegion := providers.DefaultAWSRegion
	if region != nil && len(*region) > 0 {
		defaultRegion = *region
	}

	fn := make([]providers.CredsProviderOptionsFunc, 0)
	fn = append(fn, providers.WithRegion(defaultRegion))

	if profile != nil && len(*profile) > 0 {

		fn = append(fn, providers.WithConfigProfile(*profile))

		if credsFile != nil && len(*credsFile) > 0 {
			fn = append(fn, providers.WithCredentialsFile(*credsFile))
		}

		if configFile != nil && len(*configFile) > 0 {
			fn = append(fn, providers.WithConfigFile(*configFile))
		}

		_, _ = providers.NewSharedConfigCredsProvider(context.Background(), providers.DefaultCredsProviderName, fn...)
		return
	} else if isEnv != nil && *isEnv {

		if accessKeyId != nil && len(*accessKeyId) > 0 {
			fn = append(fn, providers.WithAccessKeyIdFrom(*accessKeyId))
		}
		if secretAccessKey != nil && len(*secretAccessKey) > 0 {
			fn = append(fn, providers.WithSecretAccessKeyFrom(*secretAccessKey))
		}
		if sessionToken != nil && len(*sessionToken) > 0 {
			fn = append(fn, providers.WithSessionTokenFrom(*sessionToken))
		}

		_, _ = providers.NewEnvironmentCredsProvider(context.Background(), providers.DefaultCredsProviderName, fn...)
		return
	}

	_, _ = providers.NewDefaultCredsProvider(context.Background(), providers.DefaultCredsProviderName, fn...)
}
