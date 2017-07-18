package fetch

import (
	"context"
	"fmt"
	"sync"

	"github.com/wallix/awless/graph"
)

type Fetcher interface {
	Fetch(context.Context) (*graph.Graph, error)
	FetchByType(context.Context, string) (*graph.Graph, error)
}

type FetchResult struct {
	ResourceType string
	Err          error
	Resources    []*graph.Resource
	Objects      []interface{}
}

type Func func(context.Context) ([]*graph.Resource, []interface{}, error)

type Funcs map[string]Func

type fetcher struct {
	fetchFuncs    map[string]Func
	resourceTypes []string
}

func NewFetcher(funcs Funcs) *fetcher {
	ftr := &fetcher{
		fetchFuncs: make(Funcs),
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
		go func(t string) {
			f.fetchResource(ctx, t, results)
			wg.Done()
		}(resType)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	gph := graph.NewGraph()
	for res := range results {
		if err := res.Err; err != nil {
			return gph, err
		}
		gph.AddResource(res.Resources...)
	}

	return gph, nil
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
	resources := make([]*graph.Resource, 0)
	objects := make([]interface{}, 0)

	fn, ok := f.fetchFuncs[resourceType]
	if ok {
		resources, objects, err = fn(ctx)
	} else {
		err = fmt.Errorf("no fetch func defined for resource type '%s'", resourceType)
	}

	results <- FetchResult{
		ResourceType: resourceType,
		Err:          err,
		Resources:    resources,
		Objects:      objects,
	}
}
