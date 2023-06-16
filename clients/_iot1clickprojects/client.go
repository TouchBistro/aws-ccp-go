// AUTO-GENERATED CODE - DO NOT EDIT
// See instructions under /codegen/README.md
// GENERATED ON 2023-06-16 18:24:12

// Package _iot1clickprojects provides AWS client management functions for the iot1clickprojects
// AWS service.
//
// The Client() is a wrapper on iot1clickprojects.NewFromConfig(), which creates & caches
// the client.
//
// The Delete() clears the cached client.
//
package _iot1clickprojects

import (
	"sync"

	"github.com/TouchBistro/aws-ccp-go/providers"
	"github.com/aws/aws-sdk-go-v2/service/iot1clickprojects"
)

var cmap sync.Map

// Client builds or returns the singleton iot1clickprojects client for the supplied provider
// If functional options are supplied, they are passed as-is to the underlying NewFromConfig(...)
// for the corresponding client
func Client(provider providers.CredsProvider, optFns ...func(*iot1clickprojects.Options)) (*iot1clickprojects.Client, error) {

	if provider == nil {
		return nil, providers.ErrNilProvider
	}
	if _, ok := cmap.Load(provider.Key()); !ok {
		client := iot1clickprojects.NewFromConfig(provider.Config(), optFns...)
		cmap.Store(provider.Key(), client)
	}
	client, _ := cmap.Load(provider.Key())
	return client.(*iot1clickprojects.Client), nil
}

// Delete removes the cached iot1clickprojects client for the supplied provider; This foreces the subsequent
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
