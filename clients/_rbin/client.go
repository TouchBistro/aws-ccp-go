// AUTO-GENERATED CODE - DO NOT EDIT
// See instructions under /codegen/README.md
// GENERATED ON 2024-06-25 08:17:34

// Package _rbin provides AWS client management functions for the rbin
// AWS service.
//
// The Client() is a wrapper on rbin.NewFromConfig(), which creates & caches
// the client.
//
// The Delete() clears the cached client.
package _rbin

import (
	"sync"

	"github.com/TouchBistro/aws-ccp-go/providers"
	"github.com/aws/aws-sdk-go-v2/service/rbin"
)

var cmap sync.Map

// Client builds or returns the singleton rbin client for the supplied provider
// If functional options are supplied, they are passed as-is to the underlying NewFromConfig(...)
// for the corresponding client
func Client(provider providers.CredsProvider, optFns ...func(*rbin.Options)) (*rbin.Client, error) {

	if provider == nil {
		return nil, providers.ErrNilProvider
	}
	if _, ok := cmap.Load(provider.Key()); !ok {
		client := rbin.NewFromConfig(provider.Config(), optFns...)
		cmap.Store(provider.Key(), client)
	}
	client, _ := cmap.Load(provider.Key())
	return client.(*rbin.Client), nil
}

// Must wraps the _rbin.Client( ) function & panics if a non-nil error is returned.
func Must(provider providers.CredsProvider, optFns ...func(*rbin.Options)) *rbin.Client {

	client, err := Client(provider, optFns...)
	if err != nil {
		panic(err)
	}
	return client
}

// Delete removes the cached rbin client for the supplied provider; This foreces the subsequent
// calls to Client() for the same provider to recreate & return a new instnce.
func Delete(provider providers.CredsProvider) error {

	if provider == nil {
		return providers.ErrNilProvider
	}
	if _, ok := cmap.Load(provider.Key()); ok {
		cmap.Delete(provider.Key())
	}
	return nil
}

// Refresh discards the cached rbin client if it exists, builds & returns a new singleton instance
func Refresh(provider providers.CredsProvider, optFns ...func(*rbin.Options)) (*rbin.Client, error) {

	err := Delete(provider)
	if err != nil {
		return nil, err
	}
	return Client(provider, optFns...)
}
