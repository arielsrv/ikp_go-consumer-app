package config

import (
	"fmt"
	"regexp"
	"sync"
	"time"

	"github.com/arielsrv/go-archaius"
	"github.com/arielsrv/ikp_go-restclient/rest"
	"github.com/ugurcsen/gods-generic/lists/arraylist"
)

const (
	DefaultPoolSize      = 10
	DefaultPoolTimeout   = 200
	DefaultSocketTimeout = 500
)

const (
	RestPoolPattern   = `rest\.pool\.([-_\w]+)\..+`
	RestClientPattern = `rest\.client\.([-_\w]+)\..+`
)

var restPoolPattern = regexp.MustCompile(RestPoolPattern)
var restClientPattern = regexp.MustCompile(RestClientPattern)
var restClientFactory RESTClientFactory

var instance sync.Once

func ProvideRestClients() *RESTClientFactory {
	instance.Do(func() {
		restPoolFactory := &RESTPoolFactory{builders: map[string]*rest.RequestBuilder{}}
		poolNames := getNamesInKeys(restPoolPattern)
		for _, name := range poolNames.Values() {
			restPool := &rest.RequestBuilder{
				Timeout: time.Millisecond *
					time.Duration(TryInt(fmt.Sprintf("rest.pool.%s.pool.timeout", name), DefaultPoolTimeout)),
				ConnectTimeout: time.Millisecond *
					time.Duration(TryInt(fmt.Sprintf("rest.pool.%s.pool.connection-timeout", name), DefaultSocketTimeout)),
				CustomPool: &rest.CustomPool{
					MaxIdleConnsPerHost: TryInt(fmt.Sprintf("rest.pool.%s.pool.size", name), DefaultPoolSize),
				},
			}
			restPoolFactory.add(restPool)
			restPoolFactory.register(name, restPool)
		}

		restClientFactory = RESTClientFactory{clients: map[string]*rest.RequestBuilder{}}
		clientNames := getNamesInKeys(restClientPattern)
		for _, name := range clientNames.Values() {
			poolName := String(fmt.Sprintf("rest.client.%s.pool", name))
			pool := restPoolFactory.getPool(poolName)
			restClientFactory.register(name, pool)
		}
	})

	return &restClientFactory
}

type RESTPoolFactory struct {
	restPools arraylist.List[*rest.RequestBuilder]
	builders  map[string]*rest.RequestBuilder
}

func (r *RESTPoolFactory) add(rb *rest.RequestBuilder) {
	r.restPools.Add(rb)
}

func (r *RESTPoolFactory) register(name string, rb *rest.RequestBuilder) {
	r.builders[name] = rb
}

func (r *RESTPoolFactory) getPool(name string) *rest.RequestBuilder {
	return r.builders[name]
}

type RESTClientFactory struct {
	clients map[string]*rest.RequestBuilder
}

func (r *RESTClientFactory) register(name string, restPool *rest.RequestBuilder) {
	r.clients[name] = restPool
}

func (r *RESTClientFactory) Get(name string) *rest.RequestBuilder {
	return r.clients[name]
}

func (r *RESTClientFactory) GetClients() map[string]*rest.RequestBuilder {
	return r.clients
}

func getNamesInKeys(regex *regexp.Regexp) arraylist.List[string] {
	names := arraylist.List[string]{}

	configs := archaius.GetConfigs()
	for key := range configs {
		match := regex.FindStringSubmatch(key)
		for i := range regex.SubexpNames() {
			if i > 0 && i <= len(match) {
				group := match[1]
				if !names.Contains(group) {
					names.Add(group)
				}
			}
		}
	}
	return names
}
