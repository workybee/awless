/*
Copyright 2017 WALLIX

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package graph

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"io/ioutil"

	"github.com/wallix/awless/cloud/rdf"
	tstore "github.com/wallix/triplestore"
)

type Visitor interface {
	Visit(*Graph) error
}

type visitEachFunc func(res *Resource, depth int) error

func VisitorCollectFunc(collect *[]*Resource) visitEachFunc {
	return func(res *Resource, depth int) error {
		*collect = append(*collect, res)
		return nil
	}
}

func VisitorPrependFunc(collect *[]*Resource) visitEachFunc {
	return func(res *Resource, depth int) error {
		*collect = append([]*Resource{res}, *collect...)
		return nil
	}
}

type ParentsVisitor struct {
	From        *Resource
	Each        visitEachFunc
	IncludeFrom bool
}

func (v *ParentsVisitor) Visit(g *Graph) error {
	startNode, foreach, err := prepareRDFVisit(g, v.From, v.Each, v.IncludeFrom)
	if err != nil {
		return err
	}

	return visitBottomUp(g.store.Snapshot(), startNode, foreach)
}

type ChildrenVisitor struct {
	From        *Resource
	Each        visitEachFunc
	IncludeFrom bool
}

func (v *ChildrenVisitor) Visit(g *Graph) error {
	startNode, foreach, err := prepareRDFVisit(g, v.From, v.Each, v.IncludeFrom)
	if err != nil {
		return err
	}
	return visitTopDown(g.store.Snapshot(), startNode, foreach)
}

type SiblingsVisitor struct {
	From        *Resource
	Each        visitEachFunc
	IncludeFrom bool
}

func (v *SiblingsVisitor) Visit(g *Graph) error {
	startNode, foreach, err := prepareRDFVisit(g, v.From, v.Each, v.IncludeFrom)
	if err != nil {
		return err
	}

	return visitSiblings(g.store.Snapshot(), startNode, foreach)
}

func prepareRDFVisit(g *Graph, root *Resource, each visitEachFunc, includeRoot bool) (string, func(g tstore.RDFGraph, n string, i int) error, error) {
	rootNode := root.Id()

	foreach := func(rdfG tstore.RDFGraph, n string, i int) error {
		rT, err := resolveResourceType(rdfG, n)
		if err != nil {
			return err
		}
		res, err := g.GetResource(rT, n)
		if err != nil {
			return err
		}
		if includeRoot || !root.Same(res) {
			if err := each(res, i); err != nil {
				return err
			}
		}
		return nil
	}
	return rootNode, foreach, nil
}

func visitTopDown(snap tstore.RDFGraph, root string, each func(tstore.RDFGraph, string, int) error, distances ...int) error {
	var dist int
	if len(distances) > 0 {
		dist = distances[0]
	}

	if err := each(snap, root, dist); err != nil {
		return err
	}

	triples := snap.WithSubjPred(root, rdf.ParentOf)

	var childs []string
	for _, tri := range triples {
		n, ok := tri.Object().Resource()
		if !ok {
			return fmt.Errorf("object is not a resource identifier")
		}
		childs = append(childs, n)
	}

	sort.Strings(childs)

	for _, child := range childs {
		visitTopDown(snap, child, each, dist+1)
	}

	return nil
}

func visitBottomUp(snap tstore.RDFGraph, startNode string, each func(tstore.RDFGraph, string, int) error, distances ...int) error {
	var dist int
	if len(distances) > 0 {
		dist = distances[0]
	}

	if err := each(snap, startNode, dist); err != nil {
		return err
	}
	triples := snap.WithPredObj(rdf.ParentOf, tstore.Resource(startNode))
	var parents []string
	for _, tri := range triples {
		parents = append(parents, tri.Subject())
	}

	sort.Strings(parents)

	for _, child := range parents {
		visitBottomUp(snap, child, each, dist+1)
	}

	return nil
}

func visitSiblings(snap tstore.RDFGraph, start string, each func(tstore.RDFGraph, string, int) error, distances ...int) error {
	triples := snap.WithPredObj(rdf.ParentOf, tstore.Resource(start))

	var parents []string
	for _, tri := range triples {
		parents = append(parents, tri.Subject())
	}

	if len(parents) == 0 {
		return each(snap, start, 0)
	}

	sort.Strings(parents)

	for _, parent := range parents {
		parentTs := snap.WithSubjPred(parent, rdf.ParentOf)

		var childs []string
		for _, parentT := range parentTs {
			child, ok := parentT.Object().Resource()
			if !ok {
				return fmt.Errorf("object is not a resource identifier")
			}
			childs = append(childs, child)
		}

		sort.Strings(childs)

		startType, err := resolveResourceType(snap, start)
		if err != nil {
			return err
		}

		for _, child := range childs {
			rt, err := resolveResourceType(snap, child)
			if err != nil {
				return err
			}
			sameType := rt == startType
			if sameType {
				if err := each(snap, child, 0); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

type TreePrinter struct {
	G          *Graph
	Root       *Resource
	R          *Resource
	Onresource func(string) string
}

func NewTreePrinter(r *Resource, g *Graph) *TreePrinter {
	return &TreePrinter{
		G: g, Onresource: func(id string) string { return id },
	}
}

func (g *Graph) ExtractLineageGraph(r *Resource) (*Resource, *Graph, error) {
	gph := NewGraph()

	var child *Resource
	var lastDepth int
	var parents []*Resource

	extractParents := func(res *Resource, depth int) error {
		parents = append([]*Resource{res}, parents...)
		if child != nil {
			gph.AddParentRelation(res, child)
		}
		if lastDepth != depth {
			lastDepth = depth
			child = res
		}
		return gph.AddResource(res)
	}
	if err := g.Accept(&ParentsVisitor{From: r, Each: extractParents}); err != nil {
		return nil, nil, err
	}

	lastDepth = 0
	parent := r

	extractChildren := func(res *Resource, depth int) error {
		if lastDepth != depth {
			parent = res
		}
		if err := gph.AddResource(res); err != nil {
			return err
		}
		if err := gph.AddParentRelation(parent, res); err != nil {
			return err
		}
		return nil
	}
	if err := g.Accept(&ChildrenVisitor{From: r, Each: extractChildren, IncludeFrom: true}); err != nil {
		return nil, nil, err
	}

	root := r
	if len(parents) > 0 {
		gph.AddParentRelation(parents[len(parents)-1], r)
		root = parents[0]
	}

	err := ioutil.WriteFile("./children.triples", []byte(gph.MustMarshal()), 0600)
	if err != nil {
		return nil, nil, err
	}

	return root, gph, nil
}

func (tp *TreePrinter) Print(w io.Writer) error {
	var childrenW bytes.Buffer
	var hasChildren bool
	printWithTabs := func(r *Resource, distance int) error {
		var tabs bytes.Buffer
		tabs.WriteString(strings.Repeat(" ", distance))
		for i := 0; i < distance; i++ {
			tabs.WriteByte('\t')
		}

		display := r.String()
		if r.Same(tp.R) {
			display = tp.Onresource(tp.R.String())
		} else {
			hasChildren = true
		}
		fmt.Fprintf(&childrenW, "%s└── %s\n", tabs.String(), display)
		return nil
	}
	if err := tp.G.Accept(&ChildrenVisitor{From: tp.Root, Each: printWithTabs, IncludeFrom: true}); err != nil {
		return err
	}

	if hasChildren {
		fmt.Fprintf(w, childrenW.String())
	}

	return nil
}
