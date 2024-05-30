// AUTO-GENERATED CODE - DO NOT EDIT
// See instructions under /codegen/README.md
// GENERATED ON 2024-05-30 07:41:01

// Package _iotdeviceadvisor provides AWS client management functions for the iotdeviceadvisor
// AWS service.
//
// The Client() is a wrapper on iotdeviceadvisor.NewFromConfig(), which creates & caches
// the client.
//
// The Delete() clears the cached client.
package _iotdeviceadvisor

import (
	"sync"

	"github.com/TouchBistro/aws-ccp-go/providers"
	"github.com/aws/aws-sdk-go-v2/service/iotdeviceadvisor"
)

var cmap sync.Map

// Client builds or returns the singleton iotdeviceadvisor client for the supplied provider
// If functional options are supplied, they are passed as-is to the underlying NewFromConfig(...)
// for the corresponding client
func Client(provider providers.CredsProvider, optFns ...func(*iotdeviceadvisor.Options)) (*iotdeviceadvisor.Client, error) {

	if provider == nil {
		return nil, providers.ErrNilProvider
	}
	if _, ok := cmap.Load(provider.Key()); !ok {
		client := iotdeviceadvisor.NewFromConfig(provider.Config(), optFns...)
		cmap.Store(provider.Key(), client)
	}
	client, _ := cmap.Load(provider.Key())
	return client.(*iotdeviceadvisor.Client), nil
}

// Must wraps the _iotdeviceadvisor.Client( ) function & panics if a non-nil error is returned.
func Must(provider providers.CredsProvider, optFns ...func(*iotdeviceadvisor.Options)) *iotdeviceadvisor.Client {

	client, err := Client(provider, optFns...)
	if err != nil {
		panic(err)
	}
	return client
}

// Delete removes the cached iotdeviceadvisor client for the supplied provider; This foreces the subsequent
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

// Refresh discards the cached iotdeviceadvisor client if it exists, builds & returns a new singleton instance
func Refresh(provider providers.CredsProvider, optFns ...func(*iotdeviceadvisor.Options)) (*iotdeviceadvisor.Client, error) {

	err := Delete(provider)
	if err != nil {
		return nil, err
	}
	return Client(provider, optFns...)
}
