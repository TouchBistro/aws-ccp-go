package providers

import (
	"sync"
)

type providersMap struct {
	providers map[string]CredsProvider
	mutex     sync.Mutex
}

var pmap *providersMap

// initializes the providersMap struct & it's inner contents
func init() {
	pmap = &providersMap{
		providers: make(map[string]CredsProvider),
	}
}

// get returns the provider for the supplied key, else returns a false
func (p *providersMap) get(key string) (CredsProvider, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	provider, ok := p.providers[key]
	return provider, ok
}

// put adds the supplied key, value in the providersMap
func (p *providersMap) put(key string, val CredsProvider) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.providers[key] = val
}

// Get returns the named CredsProvider
func Get(providerName string) (CredsProvider, error) {
	if p, ok := pmap.get(providerName); ok {
		return p, nil
	} else {
		return p, ErrUnknownProvider
	}
}

// MustGet returns the named CredsProvider, or panics if an error occurs
func MustGet(providerName string) CredsProvider {
	p, err := Get(providerName)
	if err != nil {
		panic(err)
	}
	return p
}

// Default returns the 'default' provider
func Default() (CredsProvider, error) {
	return Get(DefaultCredsProviderName)
}

// Clone an existing provider for supplied providerName as the cloneName and return it
// if the provider does not exist, it returns an error
func Clone(providerName, cloneName string) (CredsProvider, error) {
	p, err := Get(providerName)
	if err != nil {
		return nil, err
	}
	pmap.put(cloneName, p)
	return p, nil
}

// MustClone clones the providerName to cloneName, or panics if an error occurs
func MustClone(providerName, cloneName string) CredsProvider {
	p, err := Clone(providerName, cloneName)
	if err != nil {
		panic(err)
	}
	return p
}
