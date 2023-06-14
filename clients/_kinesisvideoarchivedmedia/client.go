//AUTO-GENERATED CODE - DO NOT EDIT
//See instructions under /codegen/README.md
//GENERATED ON 2023-06-13 21:44:51

// Package _kinesisvideoarchivedmedia provides AWS client management functions for the kinesisvideoarchivedmedia
// AWS service.
//
// The NewClient() is a wrapper on kinesisvideoarchivedmedia.NewFromConfig(), which creates & caches
// the client.
//
// The DeleteClient() clears the cached client.
//
package _kinesisvideoarchivedmedia

import (
	"github.com/TouchBistro/aws-ccp-go/providers"
	"github.com/aws/aws-sdk-go-v2/service/kinesisvideoarchivedmedia"
	"sync"
)

var cmap sync.Map

//NewClient builds or returns the singleton kinesisvideoarchivedmedia client for the supplied provider
//If functional options are supplied, they are passed as-is to the underlying NewFromConfig(...)
//for the corresponding client
func NewClient(provider providers.CredsProvider, optFns ...func(*kinesisvideoarchivedmedia.Options)) (*kinesisvideoarchivedmedia.Client, error) {

	if provider == nil {
		return nil, providers.ErrNilProvider
	}
	if _, ok := cmap.Load(provider.Key()); !ok {
		client := kinesisvideoarchivedmedia.NewFromConfig(provider.Config(), optFns...)
		cmap.Store(provider.Key(), client)
	}
	client, _ := cmap.Load(provider.Key())
	return client.(*kinesisvideoarchivedmedia.Client), nil
}

//DeleteClient deletes the cached kinesisvideoarchivedmedia client for the supplied provider; This foreces the subsequent
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
