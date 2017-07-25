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

package main

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/wallix/awless/gen/aws"
)

func generateServicesFuncs() {
	templ, err := template.New("funcs").Funcs(template.FuncMap{
		"Title":          strings.Title,
		"ToUpper":        strings.ToUpper,
		"Join":           strings.Join,
		"ApiToInterface": aws.ApiToInterface,
	}).Parse(servicesTempl)

	if err != nil {
		panic(err)
	}

	var buff bytes.Buffer
	err = templ.Execute(&buff, aws.FetchersDefs)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(filepath.Join(SERVICES_DIR, "gen_services.go"), buff.Bytes(), 0666); err != nil {
		panic(err)
	}
}

const servicesTempl = `// Auto generated implementation for the AWS cloud service

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

package awsservices

// DO NOT EDIT - This file was automatically generated with go generate

import (
  "fmt"
	"sync"

  awssdk "github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/session"
  {{- range $index, $service := . }}
  {{- range $, $api := $service.Api }}
  "github.com/aws/aws-sdk-go/service/{{ $api }}"
  "github.com/aws/aws-sdk-go/service/{{ $api }}/{{ $api }}iface"
  {{- end }}
  {{- end }}
	"github.com/wallix/awless/cloud"
	"github.com/wallix/awless/config"
    "github.com/wallix/awless/graph"
	"github.com/wallix/awless/logger"
	"github.com/wallix/awless/template/driver"
	"github.com/wallix/awless/aws/driver"
	"github.com/wallix/awless/fetch"
	"github.com/wallix/awless/aws/fetch"
)

const accessDenied = "Access Denied"

var ServiceNames = []string{
	{{- range $index, $service := . }}
  "{{ $service.Name }}",
  {{- end }}
}

var ResourceTypes = []string {
{{- range $index, $service := . }}
    {{- range $idx, $fetcher := $service.Fetchers }}
      "{{ $fetcher.ResourceType }}",
    {{- end }}
{{- end }}
}

var ServicePerAPI = map[string]string {
{{- range $index, $service := . }}
{{- range $, $api := $service.Api }}
  "{{ $api }}": "{{ $service.Name }}",
{{- end }}
{{- end }}
}

var ServicePerResourceType = map[string]string {
{{- range $index, $service := . }}
  {{- range $idx, $fetcher := $service.Fetchers }}
  "{{ $fetcher.ResourceType }}": "{{ $service.Name }}",
  {{- end }}
{{- end }}
}

var APIPerResourceType = map[string]string {
{{- range $index, $service := . }}
  {{- range $idx, $fetcher := $service.Fetchers }}
  "{{ $fetcher.ResourceType }}": "{{ $fetcher.Api }}",
  {{- end }}
{{- end }}
}

var GlobalServices = []string{
{{- range $index, $service := . }}
    {{- if $service.Global }}
      "{{ $service.Name }}",
    {{- end }}
{{- end }}
}

{{ range $index, $service := . }}
type {{ Title $service.Name }} struct {
	fetcher fetch.Fetcher
  region string
	config config
	log *logger.Logger
	{{- range $, $api := $service.Api }}
		{{ $api }}iface.{{ ApiToInterface $api }}
	{{- end }}
}

func New{{ Title $service.Name }}(sess *session.Session, awsconf config, log *logger.Logger) cloud.Service {
  {{- if $service.Global }}
	region := "global"
	{{- else}}
	region := awssdk.StringValue(sess.Config.Region)
	{{- end}}
	return &{{ Title $service.Name }}{ 
	{{- range $, $api := $service.Api }}
		{{ApiToInterface $api }}: {{ $api }}.New(sess),
	{{- end }}
		fetcher: fetch.NewFetcher(awsfetch.Build{{ Title $service.Name }}FetchFuncs(
			&awsfetch.Config{
				Sess: sess,
				Extra: awsconf,
				Log: log,
			},
		)),
		config: awsconf,
		region: region,
		log: log,
  }
}

func (s *{{ Title $service.Name }}) Name() string {
  return "{{ $service.Name }}"
}

func (s *{{ Title $service.Name }}) Region() string {
  return s.region
}

func (s *{{ Title $service.Name }}) Drivers() []driver.Driver {
  return []driver.Driver{ 
		{{- range $, $api := $service.Api }}
		awsdriver.New{{ Title $api }}Driver(s.{{ ApiToInterface $api }}),
		{{- end }}
	}
}

func (s *{{ Title $service.Name }}) ResourceTypes() []string {
	return []string{
	{{- range $index, $fetcher := $service.Fetchers }}
		"{{ $fetcher.ResourceType }}",
	{{- end }}
	}
}

func (s *{{ Title $service.Name }}) FetchResources() (*graph.Graph, error) {
	if s.IsSyncDisabled() {
		return graph.NewGraph(), nil
	}

	allErrors := new(fetch.Error)

	ctx := context.WithValue(context.Background(), "region", s.region)
  gph, err := s.fetcher.Fetch(ctx)
	defer s.fetcher.Reset()
	
	for _, e := range *fetch.WrapError(err) {
		switch ee := e.(type) {
		case awserr.RequestFailure:
			switch ee.Message() {
			case accessDenied:
				allErrors.Add(cloud.ErrFetchAccessDenied)
			default:
				allErrors.Add(ee)
			}
		case nil:
			continue
		default:
			allErrors.Add(ee)
		}
	}

	if err := gph.AddResource(graph.InitResource(cloud.Region, s.region)); err != nil {
		return gph, err
	}

	errc := make(chan error)
	var wg sync.WaitGroup

	{{- range $index, $fetcher := $service.Fetchers }}
	if s.config.getBool("aws.{{ $service.Name }}.{{ $fetcher.ResourceType }}.sync", true) {
		list, ok := s.fetcher.Get("{{ $fetcher.ResourceType }}_objects").([]*{{ $fetcher.AWSType }})
		if !ok {
			return gph, errors.New("cannot cast to '[]*{{ $fetcher.AWSType }}' type from fetch context")
		}
		for _, r := range list {
			for _, fn := range addParentsFns["{{ $fetcher.ResourceType }}"] {
				wg.Add(1)
				go func(f addParentFn, res *{{ $fetcher.AWSType }}) {
					defer wg.Done()
					err := f(gph, res)
					if err != nil {
						errc <- err
						return
					}
				}(fn, r)
			}
		}
	}
  {{- end }}

	go func() {
		wg.Wait()
		close(errc)
	}()

	for err := range errc {
		if err != nil {
				allErrors.Add(err)
		}
	}

	if allErrors.Any() {
		return gph, allErrors
	}

	return gph, nil
}

func (s *{{ Title $service.Name }}) FetchByType(t string) (*graph.Graph, error) {
	defer s.fetcher.Reset()
	ctx := context.WithValue(context.Background(), "region", s.region)
  return s.fetcher.FetchByType(ctx, t)
}

func (s *{{ Title $service.Name }}) IsSyncDisabled() bool {
	return !s.config.getBool("aws.{{ $service.Name }}.sync", true)
}

{{ end }}`
