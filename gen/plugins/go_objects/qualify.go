// Copyright (c) 2023 under the MIT license per gql-rapid-gen/LICENSE.MD

package go_objects

import "github.com/mjdrgn/gql-rapid-gen/parser"

func (p *Plugin) Qualify(schema *parser.Schema) bool {
	return len(schema.InputObjects) > 0 || len(schema.Objects) > 0
}
