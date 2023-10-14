// Copyright (c) 2023 under the MIT license per gql-rapid-gen/LICENSE.MD

package parser

import (
	"fmt"
	"github.com/mjdrgn/gql-rapid-gen/util"
	"github.com/vektah/gqlparser/v2/ast"
	"strings"
)

type ParsedField struct {
	Name        string
	Directives  map[string][]*ParsedDirective
	Description string
	Arguments   map[string]*ParsedArgumentDef
	Type        *FieldType
}

func (pf *ParsedField) Validate() error {
	if pf.Name == "" {
		return fmt.Errorf("name is required")
	}
	for k, v := range pf.Directives {
		if v == nil {
			return fmt.Errorf("directive '%s' initialised but no values", k)
		}
		for _, d := range v {
			if d.Name != k {
				return fmt.Errorf("mismatched directives '%s' and '%s'", k, d.Name)
			}
			if err := d.Validate(); err != nil {
				return fmt.Errorf("directive '%s' failed validate: %w", d.Name, err)
			}
		}
	}
	for k, v := range pf.Arguments {
		if v.Name != k {
			return fmt.Errorf("mismatched argdefs '%s' and '%s'", k, v.Name)
		}
		if err := v.Validate(); err != nil {
			return fmt.Errorf("argdef '%s' failed validate: %w", k, err)
		}
	}
	if err := pf.Type.Validate(); err != nil {
		return fmt.Errorf("type failed validate: %w", err)
	}
	return nil
}

func (pf *ParsedField) NameTitle() string {
	return util.TitleCase(pf.Name)
}

func (pf *ParsedField) NameCamel() string {
	return util.CamelCase(pf.Name)
}

func (pf *ParsedField) NameUnder() string {
	return util.UnderCase(pf.Name)
}

func (pf *ParsedField) NameDash() string {
	return util.DashCase(pf.Name)
}

func (pf *ParsedField) Label() string {
	return util.LabelCase(pf.Name)
}

func (pf *ParsedField) HasDirective(key string) bool {
	_, ok := pf.Directives[key]
	return ok
}

func (pf *ParsedField) SingleDirective(key string) *ParsedDirective {
	v, ok := pf.Directives[key]
	if !ok {
		return nil
	}
	if len(v) == 0 {
		panic(fmt.Errorf("'%s' had '%s' directive found but len 0", pf.Name, key))
	}
	if len(v) > 1 {
		panic(fmt.Errorf("field '%s' had duplicate '%s' directives", pf.Name, key))
	}
	return v[0]
}

func (pf *ParsedField) GoStructTag() string {
	entries := make([]string, 0, 16)
	if !pf.Type.Required {
		entries = append(
			entries,
			`json:"`+pf.Name+`,omitempty"`,
			`dynamodbav:"`+pf.Name+`,omitempty"`,
		)
	} else if pf.Type.Collection {
		entries = append(
			entries,
			`json:"`+pf.Name+`,omitempty"`,
			`dynamodbav:"`+pf.Name+`,omitempty"`,
		)
	} else if pf.Type.Required {
		entries = append(
			entries,
			`json:"`+pf.Name+`"`,
			`dynamodbav:"`+pf.Name+`"`,
		)
	}
	if len(entries) > 0 {
		return "`" + strings.Join(entries, " ") + "`"
	} else {
		return ""
	}
}

// NormaliseVTL generates VTL (Velocity Template Language, for AppSync) to normalise a value as per the @normalise directive.
// The key parameter defines the variable name in VTL associated with this field, e.g. "$context.source.x"
func (pf *ParsedField) NormaliseVTL(key string) (ret string) {
	dir := pf.SingleDirective("normalise")
	if dir == nil {
		return key
	}

	ret = key

	if dir.ArgBool("force_lower") {
		ret = "$util.str.toLower(" + ret + ")"
	}

	if dir.ArgBool("trim") {
		ret = ret + ".trim()"
	}

	return ret
}

// NormaliseGo generates Go code to normalise a value as per the @normalise directive.
// The key parameter defines the variable name in VTL associated with this field, e.g. "object.myField"
func (pf *ParsedField) NormaliseGo(key string) (ret string) {
	dir := pf.SingleDirective("normalise")
	if dir == nil {
		return key
	}

	ret = key

	if dir.ArgBool("force_lower") {
		ret = "strings.ToLower(" + ret + ")"
	}

	if dir.ArgBool("trim") {
		ret = "strings.TrimSpace(" + ret + ")"
	}

	return ret
}

func parseField(f *ast.FieldDefinition) *ParsedField {
	return &ParsedField{
		Name:        f.Name,
		Directives:  parseDirectives(f.Directives),
		Description: f.Description,
		Arguments:   parseArgumentDefs(f.Arguments),
		Type:        ParseFieldType(f.Type),
	}
}

func parseFields(fl ast.FieldList) []*ParsedField {
	if fl == nil {
		return nil
	}
	fields := make([]*ParsedField, 0, len(fl))
	for _, f := range fl {
		fields = append(fields, parseField(f))
	}
	return fields
}
