package gen

import "github.com/mjdrgn/gql-rapid-gen/parser"

type Plugin interface {
	Name() string
	Order() int
	Qualify(schema *parser.Schema) bool
	Generate(schema *parser.Schema, output *Output) error
}