// Copyright 2025 Redpanda Data, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/redpanda-data/protoc-gen-go-mcp/pkg/gen"
	"github.com/redpanda-data/protoc-gen-go-mcp/pkg/generator"
	"google.golang.org/protobuf/compiler/protogen"
)

func main() {
	var flagSet flag.FlagSet
	packageSuffix := flagSet.String(
		"package_suffix",
		"mcp",
		"Generate files into a sub-package of the package containing the base .pb.go files using the given suffix. An empty suffix denotes to generate into the same package as the base pb.go files.",
	)
	excludeComments := flagSet.Bool(
		"exclude_comments",
		false,
		"Exclude protobuf leading comments from generated JSON schema description fields.",
	)
	naming := flagSet.String(
		"naming",
		"fully_qualified",
		`Tool naming convention: "fully_qualified" (pkg_Service_Method), "service_method" (Service/Method), or "method_only" (Method).`,
	)

	protogen.Options{
		ParamFunc: flagSet.Set,
	}.Run(func(plugin *protogen.Plugin) error {
		nameProvider, err := gen.ParseNameProvider(*naming)
		if err != nil {
			fmt.Fprintf(os.Stderr, "protoc-gen-go-mcp: %v\n", err)
			return err
		}
		for _, f := range plugin.Files {
			if !f.Generate {
				continue
			}
			generator.NewFileGenerator(f, plugin, generator.Options{
				ExcludeComments: *excludeComments,
				NameProvider:    nameProvider,
			}).Generate(*packageSuffix)
		}
		return nil
	})
}
