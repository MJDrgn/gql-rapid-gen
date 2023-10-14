// Copyright (c) 2023 under the MIT license per gql-rapid-gen/LICENSE.MD

package react_crud

import (
	"github.com/mjdrgn/gql-rapid-gen/gen"
	"github.com/mjdrgn/gql-rapid-gen/parser"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPlugin_Execute(t *testing.T) {
	plugin := &Plugin{}

	schema := &parser.Schema{
		Enums: map[string]*parser.ParsedEnum{
			"Test": {
				Name:        "Test",
				Description: "My Test Enum",
				Values: []*parser.EnumValue{
					{
						Name:        "V1",
						Description: "V1",
						Directives:  nil,
					},
					{
						Name:        "V2",
						Description: "V2",
						Directives:  nil,
					},
					{
						Name:        "ThirdValue",
						Description: "Third Value Description",
						Directives:  nil,
					},
				},
			},
		},
	}

	output := &gen.Output{}

	err := plugin.Generate(schema, output)
	require.Nil(t, err)

	files := output.GetFiles()

	require.Equal(t, 1, len(files))

	var file *gen.OutputFile
	for _, f := range files {
		file = f
		break
	}

	rendered := file.String()

	assert.True(t, len(rendered) > 100)
	assert.Contains(t, rendered, "V1")
	assert.Contains(t, rendered, "V2")
	assert.Contains(t, rendered, "ThirdValue")
	assert.Contains(t, rendered, "Test")
}
