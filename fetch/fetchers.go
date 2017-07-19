package fetch

import (
	"context"
	"fmt"
	"sync"

	"github.com/wallix/awless/graph"
)

type Fetcher interface {
	Cache
	Fetch(context.Context) (*graph.Graph, error)
	FetchByType(context.Context, string) (*graph.Graph, error)
}

type Cache interface {
	Store(key string, val interface{})
	Get(key string) interface{}
	Reset()
}

type FetchResult struct {
	ResourceType string
	Err          error
	Resources    []*graph.Resource
	Objects      interface{}
}

type Func func(context.Context, Cache) ([]*graph.Resource, interface{}, error)

type Funcs map[string]Func

type fetcher struct {
	*cache
	fetchFuncs    map[string]Func
	resourceTypes []string
}

func NewFetcher(funcs Funcs) *fetcher {
	ftr := &fetcher{
		fetchFuncs: make(Funcs),
		cache:      newCache(),
	}
	for resType, f := range funcs {
		ftr.resourceTypes = append(ftr.resourceTypes, resType)
		ftr.fetchFuncs[resType] = f
	}
	return ftr
}

func (f *fetcher) Fetch(ctx context.Context) (*graph.Graph, error) {
	results := make(chan FetchResult, len(f.resourceTypes))
	var wg sync.WaitGroup

	for _, resType := range f.resourceTypes {
		wg.Add(1)
		go func(t string, co context.Context) {
			f.fetchResource(co, t, results)
			wg.Done()
		}(resType, ctx)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	gph := graph.NewGraph()
	var err error
	for res := range results {
		if err = res.Err; err != nil {
			continue
		}
		gph.AddResource(res.Resources...)
	}

	return gph, err
}

func (f *fetcher) FetchByType(ctx context.Context, resourceType string) (*graph.Graph, error) {
	results := make(chan FetchResult)
	defer close(results)

	go f.fetchResource(ctx, resourceType, results)

	gph := graph.NewGraph()
	select {
	case res := <-results:
		if err := res.Err; err != nil {
			return gph, err
		}
		for _, r := range res.Resources {
			gph.AddResource(r)
		}
		return gph, nil
	}
}

func (f *fetcher) fetchResource(ctx context.Context, resourceType string, results chan<- FetchResult) {
	var err error
	var objects interface{}
	resources := make([]*graph.Resource, 0)

	fn, ok := f.fetchFuncs[resourceType]
	if ok {
		resources, objects, err = fn(ctx, f.cache)
	} else {
		err = fmt.Errorf("no fetch func defined for resource type '%s'", resourceType)
	}

	f.cache.Store(fmt.Sprintf("%s_objects", resourceType), objects)
	f.cache.Store(fmt.Sprintf("%s_resources", resourceType), resources)

	results <- FetchResult{
		ResourceType: resourceType,
		Err:          err,
		Resources:    resources,
		Objects:      objects,
	}
}

type cache struct {
	sync.RWMutex
	store map[string]interface{}
}

func newCache() *cache {
	return &cache{
		store: make(map[string]interface{}),
	}
}

func (c *cache) Store(key string, val interface{}) {
	c.Lock()
	c.store[key] = val
	c.Unlock()
}

func (c *cache) Get(key string) interface{} {
	c.RLock()
	defer c.RUnlock()
	return c.store[key]
}

func (c *cache) Reset() {
	c.Lock()
	c.store = make(map[string]interface{})
	c.Unlock()
}
