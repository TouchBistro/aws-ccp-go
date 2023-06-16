// AUTO-GENERATED CODE - DO NOT EDIT
// See instructions under /codegen/README.md
// GENERATED ON 2023-06-16 18:24:12

// Package _migrationhuborchestrator provides AWS client management functions for the migrationhuborchestrator
// AWS service.
//
// The Client() is a wrapper on migrationhuborchestrator.NewFromConfig(), which creates & caches
// the client.
//
// The Delete() clears the cached client.
//
package _migrationhuborchestrator

import (
	"sync"

	"github.com/TouchBistro/aws-ccp-go/providers"
	"github.com/aws/aws-sdk-go-v2/service/migrationhuborchestrator"
)

var cmap sync.Map

// Client builds or returns the singleton migrationhuborchestrator client for the supplied provider
// If functional options are supplied, they are passed as-is to the underlying NewFromConfig(...)
// for the corresponding client
func Client(provider providers.CredsProvider, optFns ...func(*migrationhuborchestrator.Options)) (*migrationhuborchestrator.Client, error) {

	if provider == nil {
		return nil, providers.ErrNilProvider
	}
	if _, ok := cmap.Load(provider.Key()); !ok {
		client := migrationhuborchestrator.NewFromConfig(provider.Config(), optFns...)
		cmap.Store(provider.Key(), client)
	}
	client, _ := cmap.Load(provider.Key())
	return client.(*migrationhuborchestrator.Client), nil
}

// Delete removes the cached migrationhuborchestrator client for the supplied provider; This foreces the subsequent
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
