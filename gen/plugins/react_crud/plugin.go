// Copyright (c) 2023 under the MIT license per gql-rapid-gen/LICENSE.MD

package react_crud

import "github.com/mjdrgn/gql-rapid-gen/gen"

type Plugin struct {
}

func (p *Plugin) Name() string {
	return "react_crud"
}

func (p *Plugin) Order() int {
	return 0
}

func init() {
	gen.RegisterPlugin("react_crud", &Plugin{})
}
