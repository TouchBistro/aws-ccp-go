package awscredsprovider

// This package supplies an init function such that when client code imports
// the root package with a blank import, a `DefaultCredsProvider` named `default`
// is initialized using the AWS Default Credentials Chain.
//
// The credentails from this `default` provider are not validated.
//
//
//	// intializes a `default` provider using AWS SDK defaults,
//	_ github.com/TouchBistro/aws-ccp-go
//
// This provider can be overridden by instantiating a new provider with the
// same  name, OR used as a base credentials provider for a `AssumeRoleCredsProvider`
//
