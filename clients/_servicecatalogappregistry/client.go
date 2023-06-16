// AUTO-GENERATED CODE - DO NOT EDIT
// See instructions under /codegen/README.md
// GENERATED ON 2023-06-16 18:24:12

// Package _servicecatalogappregistry provides AWS client management functions for the servicecatalogappregistry
// AWS service.
//
// The Client() is a wrapper on servicecatalogappregistry.NewFromConfig(), which creates & caches
// the client.
//
// The Delete() clears the cached client.
//
package _servicecatalogappregistry

import (
	"sync"

	"github.com/TouchBistro/aws-ccp-go/providers"
	"github.com/aws/aws-sdk-go-v2/service/servicecatalogappregistry"
)

var cmap sync.Map

// Client builds or returns the singleton servicecatalogappregistry client for the supplied provider
// If functional options are supplied, they are passed as-is to the underlying NewFromConfig(...)
// for the corresponding client
func Client(provider providers.CredsProvider, optFns ...func(*servicecatalogappregistry.Options)) (*servicecatalogappregistry.Client, error) {

	if provider == nil {
		return nil, providers.ErrNilProvider
	}
	if _, ok := cmap.Load(provider.Key()); !ok {
		client := servicecatalogappregistry.NewFromConfig(provider.Config(), optFns...)
		cmap.Store(provider.Key(), client)
	}
	client, _ := cmap.Load(provider.Key())
	return client.(*servicecatalogappregistry.Client), nil
}

// Delete removes the cached servicecatalogappregistry client for the supplied provider; This foreces the subsequent
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
