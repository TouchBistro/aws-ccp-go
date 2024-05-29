

[![GoDoc][1]][2] [![License: MIT][3]][4] [![CircleCI][5]][6] 

[1]: https://pkg.go.dev/badge/github.com/evalphobia/aws-ccp-gor?utm_source=godoc
[2]: https://pkg.go.dev/github.com/TouchBistro/aws-ccp-go
[3]: https://img.shields.io/badge/License-MIT-blue.svg
[4]: LICENSE
[5]: https://dl.circleci.com/status-badge/img/gh/TouchBistro/aws-ccp-go/tree/master.svg?style=svg
[6]: https://dl.circleci.com/status-badge/redirect/gh/TouchBistro/aws-ccp-go/tree/master

# AWS Config & Clients Provider for Go `(aws-ccp-go)`

A lightweight wrapper on the **AWS SDK for Go V2**, Config & Clients API.


`aws-ccp-go` is a light-weight wrapper on the [AWS SDK for Go v2](https://github.com/aws/aws-sdk-go-v2).
The goal of this module is to simplify the AWS Go SDK configuration sources. It also aims at making initialization of
clients with the most commonly used configuration options, simple & straightforward. It eliminates the need of having 
boilerpate code spread across the clients' code base, as well as making client initialization code across multiple 
projects consistent & easier.

For instance, initializing an EC2 client using the default credential sources of AWS SDK can be acheived by 
simply importing this module with a `_` (blank) identifier. As a side-effect of a blank import, a `CredsProvider`
named `default` is initialized using the AWS default credentials chain & cached for later use.

<br>

> More on AWS SDK configuration [here](https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/). 

<br>

This works for majority of use cases where AWS credentials are used from the most common sources, e.g AWS environment 
variable, or shared config & credentials files, IAM role for ECS tasks or IAM role for EC2 instance (in that order)


```go
 import (

	// intializes a `default` provider using AWS SDK defaults, with a DefaultCredsProvider
	_ github.com/TouchBistro/aws-ccp-go 

	// import the providers package to call the provider builder functions
	github.com/TouchBistro/aws-ccp-go/providers 
	// import the clients/_ec2 package to call the ec2 client initializer
	github.com/TouchBistro/aws-ccp-go/clients/_ec2 
 )

		:
		:
		// fetch the cached `default` provider
		p, err := providers.Default()
		// get an EC2 client
		client, err := _ec2.Client(p)
 
		:
		:
```

If multiple AWS credential sources are required, the `aws-ccp-go` API makes it simple to explicitly define 
the configuration without the need of constructing an `aws.Config` object in the client code using various builder 
functions supplied with the AWS SDK.

For instance, if the AWS static credentials are to be supplied using non-standard environment variables, you can use 
the `EnvironmentCredsProvider` to configure this. 

```go

   // uses static credentials from standard AWS env vars; 
   // AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN
   env0, err := providers.NewEnvironmentCredsProvider(ctx, "env0")

   // uses static credentials from custom env vars
   env1, err := providers.NewEnvironmentCredsProvider(ctx, "env1", 
			providers.WithAccessKeyId
			From("MY_ACCESS_KEY_ID"), 
				providers.WithSecretAccessKeyFrom("MY_SECRET_ACCES"), 
					providers.WithSessionTokenFrom("MY_SESSION_TOKEN"))

```

This allows the client code to set up multiple named `CredsProvider`s that use different sets to environment variables 
for fetching their separate static AWS credentials.

As shown in examples, `aws-ccp-go` uses a concept of named ***Providers*** to encapsulate the `aws.Config` 
and can later be used to initialize clients for AWS services

The module supplies client builders (helper functions) for all AWS services supported by the AWS SDK. These 
methods use an internal map to maintain & return singleton clients per provider.

The helper functions are exposed by service packages of this module under the `/clients/` path; and of the form:
`github.com/TouchBistro/aws-ccp-go/clients/_<service_name>`. These functions initialize & return AWS service 
clients using credentials encapsulated by these providers.

<br>

> These client builder functions supplied in the `_<service_name>` packages are all auto-generated. More details about 
the code generation [here](./codegen/README.md). An example here [client.go](./clients/_ec2/client.go) for reference.

<br>

A simple example of initializing a provider & client for ECS.

```go

  p1, _ := provider.NewDefaultCredsProvider(context.Background(),"p1")
  client, _ := _ecs.Client(p1)

```

Following are more details on the types of `providers` included with this module:

<br><br><br>

# AWS Configuration & Credentials
There are several ways to supply credentials to an AWS SDK Client. This is done by building a configuration (`aws.Config`) 
object with the sources of credentials & other required configuration attributes. When no specific configuration options 
are supplied, the SDK uses the default credentials chain as documented in the AWS SDK documentation.

The `aws-ccp-go` uses the **Providers** abstraction to intialize `aws.Config` for a specific source. All providers 
have a name. The `NewXXXCredsProvider(...)` functions return a new `provider` or a non-nil error if something goes wrong 
while configuring the provider. 

```go
	pr, err := providers.NewDefaultCredsProvider(context.Background(), "provider1")

```
This provider is also cached by the `aws-ccp-go` so it can be retrived in another location in the client code by 
using the `provider.Default()` or `providers.Get(string)` functions 

```go

	// returns a previously initialized provider named 
	// `provider1`, or a non-nil error
	pr, err := providers.Get("provider1") 

```

the `default` provider can be retrieved with 

```go
	// returns a provided named `default`, or a
	// non-nil error
	pr, err := providers.Default()
```

<br><br>
## `CredsProviderOptions`
<br>

The `NewXXXCredsProvider(...)` builder functions provided by the `aws-ccp-go` use the configuration parameters 
supplied using the [CredsProviderOptions](/providers/provider_creds_options.go) struct to initialize the corresponding 
type of provider. Not all options are used for each type of provider. The type alias `CredsProviderOptionsFunc` acts as 
functional options that can be passed to the `NewXXXCredsProvider` functions to set the required AWS `config.LoadOptions`

These `CredsProviderOptions` are passed down to the `config.LoadDefaultConfig(context, opts...)` function in the AWS SDK 
to load the `aws.Config`. The config is then encapsualted by the provider & can be later used to initialize AWS SDK clients 
for services. 

AWS `config.LoadDefaultConfig(...)` functions accepts functional options of type `config.LoadOptions` to customize a lot 
more than what is exposed by the `CredsProviderOptions`. So if it is required to configure additional load options, then 
those may be supplied as `CredsProviderOptions.LoadOptionFns` which is a slice of `func(*config.LoadOptions) error`.

A set of options supplied by `CredsProviderOptions` & how to apply to various types of `CredsProvider`s is shown in this 
table below.

| Function | Description | Applies To |
| :------- | :---------- | :--------- |
| `WithDefaultRegion()` | Sets `AwsDefaultRegion` constant value OR `(us-east-1)` as the AWS Region to use when building `aws.Config` for the provider.  | ALL |
| `WithRegion(string)` | Sets the supplied region as the AWS region to use for the `aws.Config` for the provider.  | ALL |
| `WithConfigLoadOptFns(optFns ...func(*config.LoadOptions) error) `| Supplies a set of functional options to set AWS SDK `config.LoadOptions`. This allows configuring the load options for the underlying `aws.Config` that are not availabe with the `CredsProviderOptions`. The `optFns` supplied here are passed directly as-is to the `config.LoadDefaultConfig()` calls to AWS SDK under the hood. These are applied after the other Load Options set with `CredsProviderOptions`|ALL|
| `ValidateProvider()`  | Turns on provider credentials validation. This results the provider to test the `aws.Config` by creating an internal `STS` client & making an `sts:GetCallerIdentity()` API call using the underlying configuration. This makes certain that the credentials are valid. *Note that this does not check any permissions, just the validity of AWS credentials obtained using the provider* | ALL |
| `WithAccessKeyIdFrom(string)`  | Specifies the environment variable name to read the AWS Access Key ID static credential from. | `EnvironmentCredsProvider` |
| `WithSecretAccessKeyFrom(string)` | Specifies the environment variable name to read the AWS Secret Access Key credential from. | `EnvironmentCredsProvider` |
| `WithSessionTokenFrom(string)`   | Specifies the environment variable name to read the AWS Session Token from. | `EnvironmentCredsProvider` |
| `WithConfigFile(string)`  | Specifies the AWS shared config file to use to fetch the credentials. If not supplied, the default value used is `~/.aws/config`. | `SharedConfigCredsProvider` |
| `WithCredentialsFile(string)` | Sspecifies the AWS shared credentials file to use to fetch the credentials If not supplied the default values used is `~/.aws/credentials`. | `SharedConfigCredsProvider` |
| `WithConfigProfile(string)`   | Specifies the AWS config *profile* to use. If not supplied, the default values used is `default` | `SharedConfigCredsProvider` |
| `WithBaseCredsProviderName(string)`  | Specifies the name of the base `CredsProvider` to use for assuming the supplied role.  | `AssumeRoleCredsProvider` |
| `WithBaseCredsProvider(CredsProvider)` |  Specifies the `CrdsProvider` instance to use for assuming the supplied role. If both this & `WithCredsProviderName(string)` are used, this option takes precendence. |`AssumeRoleCredsProvider` |
| `WithRoleArn(string)`   | Supplies the Role ARN to assume to obtain credentials. Role ARN takes precendence over the combination of Role Name & AWS Account ID. | `AssumeRoleCredsProvider` |
| `WithRoleName(string)`  |  Supplies the Role Name to assume to obtain credentials. `WithAccountId(string)` must also be supplied when Role Name is provided. | `AssumeRoleCredsProvider` |
| `WithAccountId(string)`   | Supplies the AWS Account ID that contains the Role Name to assume. The `WithRoleName(string)`options must also be supplied when Account ID is provided.  | `AssumeRoleCredsProvider` |


<br><br>

The following `CredsProvider`s are supported by the module at this time:

<br><br>
## `DefaultCredsProvider`
<br>
This is the default credentials provider which requires no custom configuration. This provider uses the AWS SDK's 
*default credentials chain* to find the credentials to use. The order in which the credentials are located 
depends on the AWS SDK implmenenation and is generally as follows:

- AWS **Environment Variables** (Static or Web Identitiy token)
- **Shared Configuration** files in user's profile directory
- **ECS IAM Role** for tasks. (For ECS service & standalone tasks)
- **EC2 IAM Role** (For applications running on an EC2 instances)

Here's a code sample on how to initialize a new `DefaultCredsProvider`

```go

	def, err := providers.NewDefaultCredsProvider(ctx, "def")

```

The only configuration option allowed for this provider type is the *AWS Region*. If not supplied, `us-east-1` is used as the
default AWS Region. 

<br>

> `us-east-1` is the default value used by all other providers as well.

<br>

### **Implicit `default` Provider**:

As discussed in the earlier section, when the root package of the module is imported with a blank identifier, a named 
`DefaultCredsProvider` called `default` is initialized for the calling client code. 

<br>

## Configuring `config.LoadOptions` further
If any config load options that are not directly supported by the `CredsProviderOptions` are needed, they can be
supplied via the `WithConfigLoadOptFns(optFns ...func(*config.LoadOptions) error)` functional option. The example below
shows how we can setup an alternate `EndpointResolver` so we can point the clients from this `DefaultCredsProvider` to a 
[LocalStack](https://localstack.cloud/solutions/cloud-emulation/) simulated cloud environment rather than the AWS API
endpoints.


```go

	localStackEndpoint := "http://localhost:4566"
	resolverWithOpts := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           localStackEndpoint,
			SigningRegion: "us-east-1",
			PartitionID:   "aws",
		}, nil
	})
	dep1, err := providers.NewDefaultCredsProvider(context.Background(), "def", providers.WithRegion("us-east-1"),
		providers.WithConfigLoadOptFns(config.WithEndpointResolverWithOptions(resolverWithOpts)), providers.ValidateProvider())

```


<br><br>

## `EnvironmentCredsProvider` 

<br>
This provider looks for static credentials supplied in the user's environment. The environment variable names to 
use for fetching the AWS Access Key ID, the Secret Access Key & the Session Token (where applicable) must be supplied 
using optiions. When the environment variable names are not supplied, the AWS SDK defaults are used. 

If the credentials are not found a non-nil `error` is returned.

```go
	env1, err := providers.NewEnvironmentCredsProvider(ctx, "env1",
			providers.WithAccessKeyIdFrom("MY_AWS_ACCESS_KEY_ID"),
				providers.WithSecretAccessKeyFrom("MY_AWS_SECRET_ACCESS_KEY"), 
					providers.WithSessionTokenFrom("MY_AWS_SESSION_TOKEN"),
						 providers.WithDefaultRegion())
```

<br><br>

### **Implicit `default` Environment Creds Provider**:

The `default` named provider can be auto-initialized with an `EnvironmentCredsProvider` by reading the runtime command-line arguments passed on to the calling code. This can be achieved by importing the `init/cmd` package of the module with a blank identifier. See example below for the command-line flags that `aws-ccp-go` looks for for initialization. 

```bash 

% your_util --env 
```

Here `your_util` is the binary created from the client code that imports `github.com/TouchBistro/aws-ccp-go/init/cmd` When `--env` flag is passed the `default` provider is initialized
using the AWS SDK standard environment variables for the access key id, secret access key and the session token. 

The following example shows how to use non-standard environment variables to read the AWS access credentials.
 
```bash 
% your_util --env \
            --access-key-id-from YOUR_KEY_ID_VAR \
			--secret-access-key-from YOUR_SECRET_KEY_VAR \
			--session-token-from YOUR_SESSION_TOKEN_VAR  
```

## `SharedConfigCredsProvider`

<br>
This provider looks for credentials supplied via the shared config & credentials files. By default, if 
not supplied, these files in the user's home directory at `~/.aws/config` & `~/.aws/credentials` respectively
are used. The default value for config profile is `default`. 

With this provider, you can supply custom locations and/or a named profie to use for credentials. If the 
supplied shared config files are not found a non-nil `error` is returned. 

```go
	shared1, err := providers.NewSharedConfigCredsProvider(ctx, "shared1", 
			 providers.WithConfigFile("/temp/.aws/config"), 
				providers.WithCredentialsFile("/temp/.aws/credentials"),   
					providers.WithConfigProfile("core"))
```

### **Implicit `default` Shared Configuration Creds Provider**:

The `default` named provider can be auto-initialized with a `SharedConfigCredsProvider` by reading the runtime command-line arguments passed on to the calling code. This can be achieved by importing the `init/cmd` package of the module with a blank identifier. See example below for the command-line flags that `aws-ccp-go` looks for for initialization. 

```bash 

% your_util --profile someProfileName 
```

Here `your_util` is the binary created from the client code that imports `github.com/TouchBistro/aws-ccp-go/init/cmd` When `--profile` flag is supplied with a value, the `default` provider is initialized using the AWS SDK shared configuration files for reading the AWS credentials. The default shared configuration files located in `~/.aws/config` & `~/.aws/credentials` are used.

The following example shows how to use non-default shared configuration files


```bash 
% your_util --profile someProfileName \
            --creds-file ~\.aws\credentials2 \
			--config-file ~\.aws\config2
```


<br><br>

## `AssumeRoleCredsProvider`

<br>
This provider finds AWS credentials in a 2-step process. First it uses the supplied base `CredsProvider` to initialize
an AWS STS client. If no default credential provider is supplied, an inline `DefaultCredsProvider` is initialized which 
uses the default credentials chain for finding the AWS credentials. In the second step, an `sts:assumeRole` API call 
is made to retrieve `stscreds` credentials by assuming the role that is supplied as the `role_arn` option (or a 
combination of role name & AWS accounnt id). The credentials from the base creds provider must have sufficient AWS 
permissions to assume the role supplied, else a non-nil error is returned during the provider.

Here's an example of the `AssumeRoleCredsProvider` with a role arn.

```go
    // using the env1 CredsProvider, assume the supplied role 
	// and use those credentials...
	r1, err := providers.NewAssumeRoleCredsProvider(ctx, "role1",
			 providers.WithBaseCredsProviderName("env1"), 
				providers.WithRoleArn("arn:aws:iam::123456789012:role/some-role")

```

or you can supply account id and role name.

```go
    // using the env1 CredsProvider, assume the supplied role in 
	// the said AWS account and use those credentials...
	r2, err := providers.NewAssumeRoleCredsProvider(ctx, "role2", 
		providers.WithBaseCredsProviderName("env1"), 
			providers.WithAccountId("123456789012"), 
				providers.WithRoleName("some-role"))

```

If no role is supplied, the `AssumeRoleCredsProvider` simply uses the base role credentials.

<br><br><br>
# AWS Client Builder Functions

Besides credentials, the `aws-ccp-go` module also supplies convenient builder functions for all AWS 
SDK supported clients. Each of these clients is in it's own package and must be imported to use the client 
builder functions. The package names are consistent with the AWS SDK service name packages but with a 
leading underscore `_`. The reason to use an underscore prefix in the name is to make is easier for these & AWS
SDK service packages to be imported together in the client code without the need of explicit aliases.

Also putting each client function in it's own package is done for optimization reasons. The AWS SDK repos has 
multiple go modules; and all of the service clients under the `/service/...` directories are their own module. 
Including all AWS Client builders in the same package here will result in indirectly importing all of the AWS 
SDK modules in client code & bloating the size of the resulting binaries unnecessarily.

For example, just constructing a default provider & EC2 client wrapper results in a resulting binary size `8 MB`. 
if the entire SDK was imported indirectly, this would go up to `68 MB` for just a few lines of client code. So 
with this design, the size of the imported libraries in the final go binary generated for the client is no bigger 
than what it would have been if the AWS SDK was directly used without the `aws-ccp-go` wrapper. The downside 
is having to import individual client packages alongside each AWS SDK service package.

All client wrapper packages are found under `/clients/` sub-directory in the repo.

For instance to create a AWS Client for the `applicationautoscaling` service, the corresponding `aws-ccp-go`
packages that need to be imported will be:

```go 

import (
	// To auto initialize the 'default` provider only
	_ "github.com/TouchBistro/aws-ccp-go"

	// To retrieve the default provider from the module
	"github.com/TouchBistro/aws-ccp-go/providers"

	// To construct applicationautoscaling service client 
	"github.com/TouchBistro/aws-ccp-go/clients/_applicationautoscaling"

	// To use the applicationauthoscaling AWS client provided data structures & functions
	"github.com/aws/aws-sdk-go-v2/client/applicationauthoscaling"

)
```
All client packages expose `4` helper functions: `Client()`, `Must()`, `Delete()` and `Refresh()`

The `Client` function returns a singleton AWS service client. It uses the supplied `providers.CredsProvdier` 
sub-type as the configuration (`aws.Config`).

The `Must()` function is a wrapper for the `Client()` function & panics if a non-nil error is returned. It 
allows convenience for initializing or passing AWS clients in the client code.

The `Delete` function clears the singleton instance for the supplied `provider` to force the module to create 
and return a new instance in the next call to `Client`

The `Refresh()` function discards the singleton client if it exists & recreates it.

The `aws-ccp-go` supports basic configurartion out of the box. If more specific configuration is required, 
functional options can be supplied to client builder methods for instance: AWS Region or Client Retry Attempts etc.  
These options are service-specific & are of the form `func(*<client>.Options)`. If options are supplied, they are 
forwarded as-is to the clients' `NewFromConfig(...)` builder functions.

```go 

  // build & return an ECS client for the `r2` provider with 
  // additioanl functional options, these are passed directly
  // to the underlying server.NewFromConfig(...) function for 
  // additional configuration
  client, err := _ecs.Client(r2,func(o *ecs.Options) {
			o.Region = "us-east-2" 
			o.RetryMaxAttempts = 4
	})


```

or simply 

```go 

  // build & return an ECS client for the `r2` provider 
  // Since Must() is used, an error while creating the Client
  // will result in a panic
  client := _ecs.Must(r2)


```

<br><br>
## Manually initializing AWS clients:

The underlying `aws.Config` configuration for a `provider` can be retrievd using the `Config()` method. 
After that, the AWS SDK API functions (`NewFromConfig(...)`) can be used to initialize the client. This can be 
useful in case a particular AWS service client is not exposed by this module, or if you just want to retrofit 
legacy code an only use this module to retrive `aws.Config`.

A code snippet to show how to do this using `glue` service client.


```go 

	var err error

	// create a shared config provider called `etl` using the etl profile 
	// this can be any CredsProvider 
	etl, err = providers.NewSharedConfigCredsProvider(
			context.Background(), "etl", providers.WithConfigProfile("etl"))
	if err != nil {
		log.Fatal(err.Error())
	}

	//initializing an AWS client for the `glue` service without using the helper functions.
	var client *glue.Client 

	// get the aws.Config for the `etl` provider & initialize a new glue client with options
	if cfg, err := etl.Config(); err != nil {
		log.Fatal(err)
	} else {
		client = glue.NewFromConfig(*cfg, func(opt *glue.Options) {
			opt.Region = "us-east-1"
		})
	}

	//using the client...
	out, err := client.ListJobs
		context.Background(), &glue.ListJobsInput{})

	if err != nil {
		log.Fatal(err.Erro())
	}

	for _, jobs := range out.JobNames {
		log.Info(jobs)
	}

```

To add support for any of the missing services' clients, a PR can be made, OR a feature request can be sent as per 
the guidelines [CONTRIBUTING.md](CONTRIBUTING.md)

## References:

### More details on Configuring the AWS SDK
https://aws.github.io/aws-sdk-go-v2/docs/configuring-sdk/

### Configuration & credential file settings
https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html