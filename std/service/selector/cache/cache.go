package cache

import (
	"fmt"
	"github.com/PKUJohnson/solar/std"
	"github.com/micro/go-micro/client/selector"
	"sync"
	"time"

	"github.com/micro/go-micro/registry"
	//"github.com/micro/go-micro/selector"
	//"github.com/PKUJohnson/solar/std"
)

/*
	Cache selector is a selector which uses the registry.Watcher to Cache service entries.
	It defaults to a TTL for 1 minute and causes a cache miss on the next request.
*/

type cacheSelector struct {
	so  selector.Options
	ttl time.Duration

	// registry cache
	sync.Mutex
	cache map[string][]*registry.Service
	ttls  map[string]time.Time

	// used to close or reload watcher
	reload chan bool
	exit   chan bool
}

var (
	DefaultTTL = time.Minute
)

func (c *cacheSelector) quit() bool {
	select {
	case <-c.exit:
		return true
	default:
		return false
	}
}

// cp copies a service. Because we're caching handing back pointers would
// create a race condition, so we do this instead
// its fast enough
func (c *cacheSelector) cp(current []*registry.Service) []*registry.Service {
	var services []*registry.Service

	for _, service := range current {
		// copy service
		s := new(registry.Service)
		*s = *service

		// copy nodes
		var nodes []*registry.Node
		for _, node := range service.Nodes {
			n := new(registry.Node)
			*n = *node
			nodes = append(nodes, n)
		}
		s.Nodes = nodes

		// copy endpoints
		var eps []*registry.Endpoint
		for _, ep := range service.Endpoints {
			e := new(registry.Endpoint)
			*e = *ep
			eps = append(eps, e)
		}
		s.Endpoints = eps

		// append service
		services = append(services, s)
	}

	return services
}

func (c *cacheSelector) del(service string) {
	delete(c.cache, service)
	delete(c.ttls, service)
}

func (c *cacheSelector) get(service string) ([]*registry.Service, error) {
	c.Lock()
	defer c.Unlock()

	// check the cache first
	services, ok := c.cache[service]
	ttl, kk := c.ttls[service]

	// got results, copy and return
	if ok && len(services) > 0 {
		// only return if its less than the ttl
		if kk && time.Since(ttl) < c.ttl {
			return c.cp(services), nil
		}
	}

	// cache miss or ttl expired
	// now ask the registry
	services, err := c.so.Registry.GetService(service)
	if err != nil {
		std.LogErrorc("cacheselector", err, fmt.Sprintf("fail to get service %s", service))
		return nil, err
	}

	// we didn't have any results so cache
	c.cache[service] = c.cp(services)
	c.ttls[service] = time.Now().Add(c.ttl)
	return services, nil
}

func (c *cacheSelector) set(service string, services []*registry.Service) {
	c.cache[service] = services
	c.ttls[service] = time.Now().Add(c.ttl)
}

func (c *cacheSelector) update(res *registry.Result) {
	if res == nil || res.Service == nil {
		return
	}

	c.Lock()
	defer c.Unlock()

	services, ok := c.cache[res.Service.Name]
	if !ok {
		// we're not going to cache anything
		// unless there was already a lookup
		return
	}

	if len(res.Service.Nodes) == 0 {
		switch res.Action {
		case "delete":
			c.del(res.Service.Name)
		}
		return
	}

	// existing service found
	var service *registry.Service
	var index int
	for i, s := range services {
		if s.Version == res.Service.Version {
			service = s
			index = i
		}
	}

	switch res.Action {
	case "create", "update":
		if service == nil {
			c.set(res.Service.Name, append(services, res.Service))
			return
		}

		// append old nodes to new service
		for _, cur := range service.Nodes {
			var seen bool
			for _, node := range res.Service.Nodes {
				if cur.Id == node.Id {
					seen = true
					break
				}
			}
			if !seen {
				res.Service.Nodes = append(res.Service.Nodes, cur)
			}
		}

		services[index] = res.Service
		c.set(res.Service.Name, services)
	case "delete":
		if service == nil {
			return
		}

		var nodes []*registry.Node

		// filter cur nodes to remove the dead one
		for _, cur := range service.Nodes {
			var seen bool
			for _, del := range res.Service.Nodes {
				if del.Id == cur.Id {
					seen = true
					break
				}
			}
			if !seen {
				nodes = append(nodes, cur)
			}
		}

		// still got nodes, save and return
		if len(nodes) > 0 {
			service.Nodes = nodes
			services[index] = service
			c.set(service.Name, services)
			return
		}

		// zero nodes left

		// only have one thing to delete
		// nuke the thing
		if len(services) == 1 {
			c.del(service.Name)
			return
		}

		// still have more than 1 service
		// check the version and keep what we know
		var srvs []*registry.Service
		for _, s := range services {
			if s.Version != service.Version {
				srvs = append(srvs, s)
			}
		}

		// save
		c.set(service.Name, srvs)
	}
}

// run starts the cache watcher loop
// it creates a new watcher if there's a problem
// reloads the watcher if Init is called
// and returns when Close is called
func (c *cacheSelector) run() {
	go c.tick()

	for {
		// exit early if already dead
		if c.quit() {
			return
		}

		// create new watcher
		w, err := c.so.Registry.Watch()
		if err != nil {
			std.LogErrorc("cacheselector", err, "fail to watch registry")
			time.Sleep(time.Second)
			continue
		}

		// watch for events
		if err := c.watch(w); err != nil {
			std.LogErrorc("cacheselector", err, "fail to update using watch results")
			continue
		}
	}
}

// check cache and expire on each tick
func (c *cacheSelector) tick() {
	t := time.NewTicker(3 * time.Minute)

	for {
		select {
		case <-t.C:
			// Clean all cache in case the DELETE command is missed.
			c.Lock()
			for service := range c.ttls {
				c.del(service)
			}
			c.Unlock()
		case <-c.exit:
			return
		}
	}
}

// watch loops the next event and calls update
// it returns if there's an error
func (c *cacheSelector) watch(w registry.Watcher) error {
	defer w.Stop()

	// manage this loop
	go func() {
		// wait for exit or reload signal
		select {
		case <-c.exit:
		case <-c.reload:
		}

		// stop the watcher
		w.Stop()
	}()

	for {
		res, err := w.Next()
		if err != nil {
			return err
		}
		c.update(res)
	}
}

func (c *cacheSelector) Init(opts ...selector.Option) error {
	for _, o := range opts {
		o(&c.so)
	}

	// reload the watcher
	go func() {
		select {
		case <-c.exit:
			return
		default:
			c.reload <- true
		}
	}()

	return nil
}

func (c *cacheSelector) Options() selector.Options {
	return c.so
}

func (c *cacheSelector) Select(service string, opts ...selector.SelectOption) (selector.Next, error) {
	sopts := selector.SelectOptions{
		Strategy: c.so.Strategy,
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	// get the service
	// try the cache first
	// if that fails go directly to the registry
	services, err := c.get(service)
	if err != nil {
		return nil, err
	}

	// apply the filters
	for _, filter := range sopts.Filters {
		services = filter(services)
	}

	// if there's nothing left, return
	if len(services) == 0 {
		return nil, selector.ErrNoneAvailable
	}

	return sopts.Strategy(services), nil
}

func (c *cacheSelector) Mark(service string, node *registry.Node, err error) {
	return
}

func (c *cacheSelector) Reset(service string) {
	return
}

// Close stops the watcher and destroys the cache
func (c *cacheSelector) Close() error {
	c.Lock()
	c.cache = make(map[string][]*registry.Service)
	c.Unlock()

	select {
	case <-c.exit:
		return nil
	default:
		close(c.exit)
	}
	return nil
}

func (c *cacheSelector) String() string {
	return "cache"
}

func NewSelector(opts ...selector.Option) selector.Selector {
	sopts := selector.Options{
		Strategy: selector.Random,
	}

	for _, opt := range opts {
		opt(&sopts)
	}

	if sopts.Registry == nil {
		sopts.Registry = registry.DefaultRegistry
	}

	ttl := DefaultTTL

	if sopts.Context != nil {
		if t, ok := sopts.Context.Value(ttlKey{}).(time.Duration); ok {
			ttl = t
		}
	}

	c := &cacheSelector{
		so:     sopts,
		ttl:    ttl,
		cache:  make(map[string][]*registry.Service),
		ttls:   make(map[string]time.Time),
		reload: make(chan bool, 1),
		exit:   make(chan bool),
	}

	go c.run()
	return c
}
