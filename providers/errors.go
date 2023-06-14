package providers

import "errors"

var ErrInvalidProviderName = errors.New("invalid or empty provider name supplied")

var ErrUnknownProvider = errors.New("unknown provider")

var ErrNilProvider = errors.New("nil provider")

var ErrInvalidAwsAccessKeyIdEnvValue = errors.New("emtpy or invalid value supplied for the AWS Access Key ID environment variable")

var ErrInvalidSecretAccessKeyEnvValue = errors.New("emtpy or invalid value supplied for the AWS Secret Access Key environment variable")

var ErrInvalidBaseProviderConfig = errors.New("no base credentials provider found")
