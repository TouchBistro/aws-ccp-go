package providers

import (
	"sync"
)

type providersMap struct {
	providers map[string]CredsProvider
	mutex     sync.Mutex
}

var pmap *providersMap

//initializes the providersMap struct & it's inner contents
func init() {
	pmap = &providersMap{
		providers: make(map[string]CredsProvider),
	}
}

//get returns the provider for the supplied key, else returns a false
func (p *providersMap) get(key string) (CredsProvider, bool) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	provider, ok := p.providers[key]
	return provider, ok
}

//put adds the supplied key, value in the providersMap
func (p *providersMap) put(key string, val CredsProvider) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.providers[key] = val
}

//contains returns a `true` boolean if a value for the key exists in the providerMap
// func (p *providersMap) contains(key string) bool {
// 	p.mutex.Lock()
// 	defer p.mutex.Unlock()
// 	_, ok := p.providers[key]
// 	return ok
// }

//Get returns the named CredsProvider
func Get(providerName string) (CredsProvider, error) {
	if p, ok := pmap.get(providerName); ok {
		return p, nil
	} else {
		return p, ErrUnknownProvider
	}
}

//Default returns the 'default' provider
func Default() (CredsProvider, error) {
	return Get(DefaultCredsProviderName)
}
