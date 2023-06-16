// AUTO-GENERATED CODE - DO NOT EDIT
// See instructions under /codegen/README.md
// GENERATED ON 2023-06-16 18:24:12

// Package _resourceexplorer2 provides AWS client management functions for the resourceexplorer2
// AWS service.
//
// The Client() is a wrapper on resourceexplorer2.NewFromConfig(), which creates & caches
// the client.
//
// The Delete() clears the cached client.
//
package _resourceexplorer2

import (
	"sync"

	"github.com/TouchBistro/aws-ccp-go/providers"
	"github.com/aws/aws-sdk-go-v2/service/resourceexplorer2"
)

var cmap sync.Map

// Client builds or returns the singleton resourceexplorer2 client for the supplied provider
// If functional options are supplied, they are passed as-is to the underlying NewFromConfig(...)
// for the corresponding client
func Client(provider providers.CredsProvider, optFns ...func(*resourceexplorer2.Options)) (*resourceexplorer2.Client, error) {

	if provider == nil {
		return nil, providers.ErrNilProvider
	}
	if _, ok := cmap.Load(provider.Key()); !ok {
		client := resourceexplorer2.NewFromConfig(provider.Config(), optFns...)
		cmap.Store(provider.Key(), client)
	}
	client, _ := cmap.Load(provider.Key())
	return client.(*resourceexplorer2.Client), nil
}

// Delete removes the cached resourceexplorer2 client for the supplied provider; This foreces the subsequent
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
