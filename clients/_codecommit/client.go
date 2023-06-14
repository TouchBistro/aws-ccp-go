//AUTO-GENERATED CODE - DO NOT EDIT
//See instructions under /codegen/README.md
//GENERATED ON 2023-06-13 21:44:51

// Package _codecommit provides AWS client management functions for the codecommit
// AWS service.
//
// The NewClient() is a wrapper on codecommit.NewFromConfig(), which creates & caches
// the client.
//
// The DeleteClient() clears the cached client.
//
package _codecommit

import (
	"github.com/TouchBistro/aws-ccp-go/providers"
	"github.com/aws/aws-sdk-go-v2/service/codecommit"
	"sync"
)

var cmap sync.Map

//NewClient builds or returns the singleton codecommit client for the supplied provider
//If functional options are supplied, they are passed as-is to the underlying NewFromConfig(...)
//for the corresponding client
func NewClient(provider providers.CredsProvider, optFns ...func(*codecommit.Options)) (*codecommit.Client, error) {

	if provider == nil {
		return nil, providers.ErrNilProvider
	}
	if _, ok := cmap.Load(provider.Key()); !ok {
		client := codecommit.NewFromConfig(provider.Config(), optFns...)
		cmap.Store(provider.Key(), client)
	}
	client, _ := cmap.Load(provider.Key())
	return client.(*codecommit.Client), nil
}

//DeleteClient deletes the cached codecommit client for the supplied provider; This foreces the subsequent
//calls to NewClient() for the same provider to recreate & return a new instnce.
func DeleteClient(provider providers.CredsProvider) error {

	if provider == nil {
		return providers.ErrNilProvider
	}
	if _, ok := cmap.Load(provider.Key()); ok {
		cmap.Delete(provider.Key())
	}
	return nil
}
