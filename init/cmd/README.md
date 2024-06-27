
# Implicit configuration of `default` provider using command-line parameters
 When aws-ccp-go clients import this piackage with a blank import, the following command line flags are parsed & used to configure the "default" creds provider.

 Importing this package is non-intrusive, as it should not affect any command-line processing for the client code. However, some third-party packages will throw an error if the flags supplied at runtime to configure aws-ccp-go are not defined & registered with those packages. So any subset of flags that the client code intends on using may also require to be configured with the client code's command-line processing library.

 # General

  `region` :  When supplied, configures the default AWS region to use. 'us-east-1' is the hard-coded default currently.

>*This flag only impacts the auto-initialized `default` provider*

## Shared Credential Files

 `profile` : When supplied, the auto-initialized default provider will be configured using a `SharedCredsProvider`
 and will use the supplied value as the profile name. This value overrides any other option on the command-line.

 ### Shared Credential Configuration Tuning:

 if `profile` flag is supplied, the following two options can be used to configure the AWS config & credential file
 location

- `creds-file` :  When supplied, the default provider uses the supplied file as the AWS Credentials file for shared config
                else default `~/.aws/credentials` is used.

- `config-file` :  When supplied,the default provider uses the supplied file as the AWS Config file for shared config
                else default `~/.aws/config` is used.

## Environment Variables

`env` : When supplied, explicitly indicates the default provider should use the EnvironmentCredsProvider & the default
        env var names to get the access_key_id, secret_access_key & session_token values. if supplied with `profile` flag
        the profile flag takes precendece & the auto-initialized default provider will use the SharedConfigCredsProvider
        instead

 ### Environment Variables Configuration Tuning:

`access-key-id-from`: When `env` flag is supplied, the valud for this flag specifies the env var name to use for fetching
                    the AWS Acces Key ID. If not supplied, AWS_ACCESS_KEY_ID is the deafult value used

`secret-access-key-from`: When `env` flag is supplied, the valud for this flag specifies the env var name to use for fetching
                        the AWS Secret Acces Key. If not supplied, AWS_SECRET_ACCESS_KEY_ID is the deafult value used

`session-token-from`: When `env` flag is supplied, the valud for this flag specifies the env var name to use for fetching
                    the AWS Session Token. If not supplied, AWS_SESSION_TOKEN is the deafult value used.


## Credential by assuming an IAM Role

`role-arn`: When this value is supplied, the `default` provider uses the `AssumeRoleCredsProvider` & the supplied role arn as the role to assume. This option works in conjunction will all others above, using the same rules, however, after the `default`
provider has been initialized with either the Environment variables or the shared credentials file; it is then used to assume the role arn supplied by this flag. If `profile` or `env` options are not supplied, then the Base credentials are taken using the **Default Credentials Chain** of the AWS SDK.