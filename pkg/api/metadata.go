/*
 * Copyright (c) 2017, MegaEase
 * All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package api

import (
	"fmt"
	"reflect"
	"sort"

	"github.com/megaease/easegress/pkg/object/httppipeline"
	"github.com/megaease/easegress/pkg/v"

	"github.com/kataras/iris"
	yaml "gopkg.in/yaml.v2"
)

const (
	// MetadataPrefix is the metadata prefix.
	MetadataPrefix = "/metadata"

	// ObjectMetadataPrefix is the object metadata prefix.
	ObjectMetadataPrefix = "/metadata/objects"

	// FilterMetaPrefix is the filter of HTTPPipeline metadata prefix.
	FilterMetaPrefix = "/metadata/objects/httppipeline/filters"
)

type (
	// FilterMeta is the metadata of filter.
	FilterMeta struct {
		Kind        string
		Results     []string
		SpecType    reflect.Type
		Description string
	}
)

var (
	filterMetaBook = map[string]*FilterMeta{}
	filterKinds    []string
)

func (s *Server) setupMetadaAPIs() {
	filterRegistry := httppipeline.GetFilterRegistry()
	for kind, f := range filterRegistry {
		filterMetaBook[kind] = &FilterMeta{
			Kind:        kind,
			Results:     f.Results(),
			SpecType:    reflect.TypeOf(f.DefaultSpec()),
			Description: f.Description(),
		}
		filterKinds = append(filterKinds, kind)
		sort.Strings(filterMetaBook[kind].Results)
	}
	sort.Strings(filterKinds)

	metadataAPIs := make([]*APIEntry, 0)
	metadataAPIs = append(metadataAPIs,
		&APIEntry{
			Path:    FilterMetaPrefix,
			Method:  "GET",
			Handler: s.listFilters,
		},
		&APIEntry{
			Path:    FilterMetaPrefix + "/{kind:string}" + "/description",
			Method:  "GET",
			Handler: s.getFilterDescription,
		},
		&APIEntry{
			Path:    FilterMetaPrefix + "/{kind:string}" + "/schema",
			Method:  "GET",
			Handler: s.getFilterSchema,
		},
		&APIEntry{
			Path:    FilterMetaPrefix + "/{kind:string}" + "/results",
			Method:  "GET",
			Handler: s.getFilterResults,
		},
	)

	s.RegisterAPIs(metadataAPIs)
}

func (s *Server) listFilters(ctx iris.Context) {
	buff, err := yaml.Marshal(filterKinds)
	if err != nil {
		panic(fmt.Errorf("marshal %#v to yaml failed: %v", filterKinds, err))
	}

	ctx.Header("Content-Type", "text/vnd.yaml")
	ctx.Write(buff)
}

func (s *Server) getFilterDescription(ctx iris.Context) {
	kind := ctx.Params().Get("kind")

	fm, exits := filterMetaBook[kind]
	if !exits {
		HandleAPIError(ctx, iris.StatusNotFound, fmt.Errorf("not found"))
		return
	}

	ctx.WriteString(fm.Description)
}

func (s *Server) getFilterSchema(ctx iris.Context) {
	kind := ctx.Params().Get("kind")

	fm, exits := filterMetaBook[kind]
	if !exits {
		HandleAPIError(ctx, iris.StatusNotFound, fmt.Errorf("not found"))
		return
	}

	buff, err := v.GetSchemaInYAML(fm.SpecType)
	if err != nil {
		panic(fmt.Errorf("get schema for %v failed: %v", fm.Kind, err))
	}

	ctx.Header("Content-Type", "text/vnd.yaml")
	ctx.Write(buff)
}

func (s *Server) getFilterResults(ctx iris.Context) {
	kind := ctx.Params().Get("kind")

	fm, exits := filterMetaBook[kind]
	if !exits {
		HandleAPIError(ctx, iris.StatusNotFound, fmt.Errorf("not found"))
		return
	}

	buff, err := yaml.Marshal(fm.Results)
	if err != nil {
		panic(fmt.Errorf("marshal %#v to yaml failed: %v", fm.Results, err))
	}

	ctx.Header("Content-Type", "text/vnd.yaml")
	ctx.Write(buff)
}
