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

func generateNewFetcherFuncs() {
	templ, err := template.New("funcs").Funcs(template.FuncMap{
		"Title":          strings.Title,
		"ToUpper":        strings.ToUpper,
		"Join":           strings.Join,
		"ApiToInterface": aws.ApiToInterface,
	}).Parse(newFetchersTempl)

	if err != nil {
		panic(err)
	}

	var buff bytes.Buffer
	err = templ.Execute(&buff, aws.FetchersDefs)
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(filepath.Join(filepath.Join(ROOT_DIR, "fetch", "aws"), "gen_fetchers.go"), buff.Bytes(), 0666); err != nil {
		panic(err)
	}
}

const newFetchersTempl = `// Auto generated implementation for the AWS cloud service

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

package awsfetch

// DO NOT EDIT - This file was automatically generated with go generate

import (
  "context"
 
  awssdk "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-sdk-go/aws/session"
  {{- range $index, $service := . }}
  {{- range $, $api := $service.Api }}
  "github.com/aws/aws-sdk-go/service/{{ $api }}"
  "github.com/aws/aws-sdk-go/service/{{ $api }}/{{ $api }}iface"
  {{- end }}
  {{- end }}
  "github.com/wallix/awless/fetch"
  "github.com/wallix/awless/graph"
  "github.com/wallix/awless/aws"
)

{{- range $index, $service := . }}
func Build{{ Title $service.Name }}FetchFuncs(sess *session.Session) fetch.Funcs {
{{- range $, $api := $service.Api }}
	{{ $api }}_api := {{ $api }}.New(sess)
	{{ $api }}_api = {{ $api }}_api // avoid not used message when api is only manual mode
{{- end }}
	
	funcs := make(map[string]fetch.Func)

	addManual{{ Title $service.Name }}FetchFuncs(sess, funcs)
	
{{- range $index, $fetcher := $service.Fetchers }}
	{{- if not $fetcher.ManualFetcher }}

	funcs["{{ $fetcher.ResourceType }}"] = func(ctx context.Context) ([]*graph.Resource, interface{}, error) {
		var resources []*graph.Resource
		var objects []*{{ $fetcher.AWSType }}
		{{- if $fetcher.Multipage }}
		var badResErr error
		err := {{ $fetcher.Api}}_api.{{ $fetcher.ApiMethod }}(&{{ $fetcher.Input }},
			func(out *{{ $fetcher.Output }}, lastPage bool) (shouldContinue bool) {
				{{- if ne $fetcher.OutputsContainers "" }}
				for _, all := range out.{{ $fetcher.OutputsContainers }} {
				{{- end }}
					for _, output := range {{ if ne $fetcher.OutputsContainers "" }}all{{ else }}out{{ end }}.{{ $fetcher.OutputsExtractor }} {
						if badResErr != nil {
							return false
						}
						objects = append(objects, output)
						var res *graph.Resource
						if res, badResErr = aws.NewResource(output); badResErr != nil {
							return false
						}
						resources = append(resources, res)
					}
				{{- if ne $fetcher.OutputsContainers "" }}
				}
				{{- end }}
				return out.{{ $fetcher.NextPageMarker }} != nil
			})
		if err != nil {
			return resources, objects, err
		}

		return resources, objects, badResErr
		{{- else }}
		
		out, err := {{ $fetcher.Api}}_api.{{ $fetcher.ApiMethod }}(&{{ $fetcher.Input }})
		if err != nil {
			return resources, objects, err
		}

		for _, output := range out.{{ $fetcher.OutputsExtractor }} {
			objects = append(objects, output)
			res, err := aws.NewResource(output)
			if err != nil {
				return resources, objects, err
			}
			resources = append(resources, res)
		}
			
		return resources, objects, nil{{ end }}
	}
{{- end }}
{{- end }}
	return funcs
}
{{- end }}`
