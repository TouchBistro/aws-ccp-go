// AUTO-GENERATED CODE - DO NOT EDIT
// See instructions under /codegen/README.md
// GENERATED ON 2023-06-16 18:24:12

// Package _lookoutmetrics provides AWS client management functions for the lookoutmetrics
// AWS service.
//
// The Client() is a wrapper on lookoutmetrics.NewFromConfig(), which creates & caches
// the client.
//
// The Delete() clears the cached client.
//
package _lookoutmetrics

import (
	"sync"

	"github.com/TouchBistro/aws-ccp-go/providers"
	"github.com/aws/aws-sdk-go-v2/service/lookoutmetrics"
)

var cmap sync.Map

// Client builds or returns the singleton lookoutmetrics client for the supplied provider
// If functional options are supplied, they are passed as-is to the underlying NewFromConfig(...)
// for the corresponding client
func Client(provider providers.CredsProvider, optFns ...func(*lookoutmetrics.Options)) (*lookoutmetrics.Client, error) {

	if provider == nil {
		return nil, providers.ErrNilProvider
	}
	if _, ok := cmap.Load(provider.Key()); !ok {
		client := lookoutmetrics.NewFromConfig(provider.Config(), optFns...)
		cmap.Store(provider.Key(), client)
	}
	client, _ := cmap.Load(provider.Key())
	return client.(*lookoutmetrics.Client), nil
}

// Delete removes the cached lookoutmetrics client for the supplied provider; This foreces the subsequent
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
