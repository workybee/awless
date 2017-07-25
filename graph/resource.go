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
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"

	"github.com/wallix/awless/cloud/rdf"
	tstore "github.com/wallix/triplestore"
)

type Resource struct {
	kind, id string

	Properties map[string]interface{}
	Relations  map[string][]*Resource
	Meta       map[string]interface{}
}

const notFoundResourceType = "notfound"

func NotFoundResource(id string) *Resource {
	return &Resource{
		id:         id,
		kind:       notFoundResourceType,
		Properties: make(map[string]interface{}),
		Meta:       make(map[string]interface{}),
		Relations:  make(map[string][]*Resource),
	}
}

func InitResource(kind, id string) *Resource {
	return &Resource{
		id:         id,
		kind:       kind,
		Properties: make(map[string]interface{}),
		Meta:       make(map[string]interface{}),
		Relations:  make(map[string][]*Resource),
	}
}

func (res *Resource) String() string {
	var identifier string
	if res == nil || (res.Id() == "" && res.Type() == "") {
		return "[none]"
	}
	if res.kind == notFoundResourceType {
		return fmt.Sprintf("%s[!not found!]", res.Id())
	}
	if name, ok := res.Properties["Name"]; ok && name != "" {
		identifier = fmt.Sprintf("@%v", name)
	} else {
		identifier = res.Id()
	}
	return fmt.Sprintf("%s[%s]", identifier, res.Type())
}

func (res *Resource) Type() string {
	return res.kind
}

func (res *Resource) Id() string {
	return res.id
}

// Compare only the id and type of the resources (no properties nor meta)
func (res *Resource) Same(other *Resource) bool {
	if res == nil && other == nil {
		return true
	}
	if res == nil || other == nil {
		return false
	}
	return res.Id() == other.Id() && res.Type() == other.Type()
}

func (res *Resource) marshalFullRDF() ([]tstore.Triple, error) {
	var triples []tstore.Triple

	cloudType := namespacedResourceType(res.Type())
	triples = append(triples, tstore.SubjPred(res.id, rdf.RdfType).Resource(cloudType))

	for key, value := range res.Meta {
		if key == "diff" {
			triples = append(triples, tstore.SubjPred(res.id, MetaPredicate).StringLiteral(fmt.Sprint(value)))
		}
	}

	for key, value := range res.Properties {
		if value == nil {
			continue
		}

		propId, err := rdf.Properties.GetRDFId(key)
		if err != nil {
			return triples, fmt.Errorf("resource %s: marshalling property: %s", res, err)
		}

		propType, err := rdf.Properties.GetDefinedBy(propId)
		if err != nil {
			return triples, fmt.Errorf("resource %s: marshalling property: %s", res, err)
		}
		dataType, err := rdf.Properties.GetDataType(propId)
		if err != nil {
			return triples, fmt.Errorf("resource %s: marshalling property: %s", res, err)
		}
		switch propType {
		case rdf.RdfsLiteral, rdf.RdfsClass:
			obj, err := marshalToRdfObject(value, propType, dataType)
			if err != nil {
				return triples, fmt.Errorf("resource %s: marshalling property '%s': %s", res, key, err)
			}
			triples = append(triples, tstore.SubjPred(res.Id(), propId).Object(obj))
		case rdf.RdfsList:
			switch dataType {
			case rdf.XsdString:
				list, ok := value.([]string)
				if !ok {
					return triples, fmt.Errorf("resource %s: marshalling property '%s': expected a string slice, got a %T", res, key, value)
				}
				for _, l := range list {
					triples = append(triples, tstore.SubjPred(res.id, propId).StringLiteral(l))
				}
			case rdf.RdfsClass:
				list, ok := value.([]string)
				if !ok {
					return triples, fmt.Errorf("resource %s: marshalling property '%s': expected a string slice, got a %T", res, key, value)
				}
				for _, l := range list {
					triples = append(triples, tstore.SubjPred(res.id, propId).Resource(l))
				}
			case rdf.NetFirewallRule:
				list, ok := value.([]*FirewallRule)
				if !ok {
					return triples, fmt.Errorf("resource %s: marshalling property '%s': expected a firewall rule slice, got a %T", res, key, value)
				}
				for _, r := range list {
					ruleId := randomRdfId()
					triples = append(triples, tstore.SubjPred(res.id, propId).Resource(ruleId))
					triples = append(triples, r.marshalToTriples(ruleId)...)
				}
			case rdf.NetRoute:
				list, ok := value.([]*Route)
				if !ok {
					return triples, fmt.Errorf("resource %s: marshalling property '%s': expected a route slice, got a %T", res, key, value)
				}
				for _, r := range list {
					routeId := randomRdfId()
					triples = append(triples, tstore.SubjPred(res.id, propId).Resource(routeId))
					triples = append(triples, r.marshalToTriples(routeId)...)
				}
			case rdf.Grant:
				list, ok := value.([]*Grant)
				if !ok {
					return triples, fmt.Errorf("resource %s: marshalling property '%s': expected a grant slice, got a %T", res, key, value)
				}
				for _, g := range list {
					grantId := randomRdfId()
					triples = append(triples, tstore.SubjPred(res.id, propId).Resource(grantId))
					triples = append(triples, g.marshalToTriples(grantId)...)
				}
			case rdf.KeyValue:
				list, ok := value.([]*KeyValue)
				if !ok {
					return triples, fmt.Errorf("resource %s: marshalling property '%s': expected a keyvalue slice, got a %T", res, key, value)
				}
				for _, kv := range list {
					keyValId := randomRdfId()
					triples = append(triples, tstore.SubjPred(res.id, propId).Resource(keyValId))
					triples = append(triples, kv.marshalToTriples(keyValId)...)
				}
			case rdf.DistributionOrigin:
				list, ok := value.([]*DistributionOrigin)
				if !ok {
					return triples, fmt.Errorf("resource %s: marshalling property '%s': expected a distribution origin slice, got a %T", res, key, value)
				}
				for _, o := range list {
					keyValId := randomRdfId()
					triples = append(triples, tstore.SubjPred(res.id, propId).Resource(keyValId))
					triples = append(triples, o.marshalToTriples(keyValId)...)
				}
			case rdf.Grant:
			default:
				return triples, fmt.Errorf("resource %s: marshalling property '%s': unexpected rdfs:DataType: %s", res, key, dataType)
			}

		default:
			return triples, fmt.Errorf("resource %s: marshalling property '%s': unexpected rdfs:isDefinedBy: %s", res, key, propType)
		}

	}
	return triples, nil
}

func marshalToRdfObject(i interface{}, definedBy, dataType string) (tstore.Object, error) {
	switch definedBy {
	case rdf.RdfsLiteral:
		return tstore.ObjectLiteral(i)
	case rdf.RdfsClass:
		return tstore.Resource(fmt.Sprint(i)), nil
	default:
		return nil, fmt.Errorf("unexpected rdfs:isDefinedBy: %s", definedBy)
	}
}

func (res *Resource) unmarshalFullRdf(gph tstore.RDFGraph) error {
	cloudType := namespacedResourceType(res.Type())
	if !gph.Contains(tstore.SubjPred(res.Id(), rdf.RdfType).Resource(cloudType)) {
		return fmt.Errorf("triple <%s><%s><%s> not found in graph", res.Id(), rdf.RdfType, cloudType)
	}
	for _, t := range gph.WithSubject(res.Id()) {
		pred := t.Predicate()
		if !rdf.Properties.IsRDFProperty(pred) || rdf.Properties.IsRDFSubProperty(pred) {
			continue
		}

		propKey, err := rdf.Properties.GetLabel(pred)
		if err != nil {
			return fmt.Errorf("unmarshalling property: label: %s", err)
		}
		propVal, err := getPropertyValue(gph, t.Object(), pred)
		if err != nil {
			return fmt.Errorf("unmarshalling property %s: val: %s", propKey, err)
		}
		if rdf.Properties.IsRDFList(pred) {
			dataType, err := rdf.Properties.GetDataType(pred)
			if err != nil {
				return fmt.Errorf("unmarshalling property: datatype: %s", err)
			}
			switch dataType {
			case rdf.RdfsClass, rdf.XsdString:
				list, ok := res.Properties[propKey].([]string)
				if !ok {
					list = []string{}
				}
				list = append(list, propVal.(string))
				res.Properties[propKey] = list
			case rdf.NetFirewallRule:
				list, ok := res.Properties[propKey].([]*FirewallRule)
				if !ok {
					list = []*FirewallRule{}
				}
				list = append(list, propVal.(*FirewallRule))
				res.Properties[propKey] = list
			case rdf.NetRoute:
				list, ok := res.Properties[propKey].([]*Route)
				if !ok {
					list = []*Route{}
				}
				list = append(list, propVal.(*Route))
				res.Properties[propKey] = list
			case rdf.Grant:
				list, ok := res.Properties[propKey].([]*Grant)
				if !ok {
					list = []*Grant{}
				}
				list = append(list, propVal.(*Grant))
				res.Properties[propKey] = list
			case rdf.KeyValue:
				list, ok := res.Properties[propKey].([]*KeyValue)
				if !ok {
					list = []*KeyValue{}
				}
				list = append(list, propVal.(*KeyValue))
				res.Properties[propKey] = list
			case rdf.DistributionOrigin:
				list, ok := res.Properties[propKey].([]*DistributionOrigin)
				if !ok {
					list = []*DistributionOrigin{}
				}
				list = append(list, propVal.(*DistributionOrigin))
				res.Properties[propKey] = list
			default:
				return fmt.Errorf("unmarshalling property: unexpected datatype %s", dataType)
			}
		} else {
			res.Properties[propKey] = propVal
		}
	}
	return nil
}

func (r *Resource) unmarshalMeta(gph tstore.RDFGraph) error {
	for _, t := range gph.WithSubjPred(r.Id(), MetaPredicate) {
		text, err := tstore.ParseString(t.Object())
		if err != nil {
			return err
		}
		r.Meta["diff"] = text
	}
	return nil
}

func namespacedResourceType(typ string) string {
	return fmt.Sprintf("%s:%s", rdf.CloudOwlNS, strings.Title(typ))
}

type Resources []*Resource

func (res Resources) Map(f func(*Resource) string) (out []string) {
	for _, r := range res {
		out = append(out, f(r))
	}
	return
}

func Subtract(one, other map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for propK, propV := range one {
		var found bool
		if otherV, ok := other[propK]; ok {
			if reflect.DeepEqual(propV, otherV) {
				found = true
			}
		}
		if !found {
			result[propK] = propV
		}
	}

	return result
}

var errTypeNotFound = errors.New("resource type not found")

func resolveResourceType(g tstore.RDFGraph, id string) (string, error) {
	typeTs := g.WithSubjPred(id, rdf.RdfType)
	switch len(typeTs) {
	case 0:
		return "", errTypeNotFound
	case 1:
		return unmarshalResourceType(typeTs[0].Object())
	default:
		return "", fmt.Errorf("cannot resolve unique type for resource '%s', got: %v", id, typeTs)
	}
}

func lowerFirstLetter(s string) string {
	a := []rune(s)
	a[0] = unicode.ToLower(a[0])
	return string(a)
}

func unmarshalResourceType(obj tstore.Object) (string, error) {
	node, ok := obj.Resource()
	if !ok {
		return "", fmt.Errorf("object is not a resource identifier, %v", obj)
	}
	return lowerFirstLetter(trimNS(node)), nil
}

func trimNS(s string) string {
	spl := strings.Split(s, ":")
	if len(spl) == 0 {
		return s
	}
	return spl[len(spl)-1]
}
