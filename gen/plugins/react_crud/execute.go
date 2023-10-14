// Copyright (c) 2023 under the MIT license per gql-rapid-gen/LICENSE.MD

package react_crud

import (
	"fmt"
	"github.com/mjdrgn/gql-rapid-gen/gen"
	"github.com/mjdrgn/gql-rapid-gen/parser"
	"path/filepath"
)

type baseData struct {
	Object  *parser.ParsedObject
	HashKey *parser.ParsedField
	SortKey *parser.ParsedField
	HasSort bool
	Fields  []field
}

type field struct {
	Def          *parser.ParsedField
	Render       string
	Input        string
	Fkey         bool
	FkeyResolver string
	FkeyField    string
}

type fieldData struct {
	Field *parser.ParsedField
	Dir   *parser.ParsedDirective
}

func (p *Plugin) Generate(schema *parser.Schema, output *gen.Output) error {

	for _, o := range schema.Objects {

		base := &baseData{
			Object: o,
		}

		if ddb := o.SingleDirective("dynamodb"); ddb != nil {
			base.HashKey = o.Field(ddb.Arg("hash_key"))
			base.SortKey = o.Field(ddb.Arg("sort_key"))
			base.HasSort = base.SortKey != nil
		}

		base.Fields = make([]field, 0, len(o.Fields))
		for _, f := range o.Fields {
			dir := f.SingleDirective("crud_type")
			renderMode := ""
			inputMode := ""

			if dir != nil && dir.HasArg("render") {
				renderMode = dir.Arg("render")
			}
			if dir != nil && dir.HasArg("input") {
				inputMode = dir.Arg("input")
			}

			if renderMode == "" {
				renderMode = "raw"
			}
			if inputMode == "" {
				inputMode = "raw"
			}

			if f.Type.Collection {
				if renderMode != "fkey" {
					renderMode = renderMode + "_collection"
				}
				if inputMode != "fkey" {
					inputMode = inputMode + "_collection"
				}
			}

			fd := fieldData{
				Field: f,
				Dir:   dir,
			}

			render, err := gen.ExecuteTemplate("plugins/react_crud/templates/render/"+renderMode+".tsx.tmpl", fd)
			if err != nil {
				return fmt.Errorf("failed rendering field render %s.%s: %w", o.Name, f.Name, err)
			}

			input, err := gen.ExecuteTemplate("plugins/react_crud/templates/input/"+inputMode+".tsx.tmpl", fd)
			if err != nil {
				return fmt.Errorf("failed rendering field input %s.%s: %w", o.Name, f.Name, err)
			}

			base.Fields = append(base.Fields, field{
				Def:          f,
				Render:       render,
				Input:        input,
				Fkey:         renderMode == "fkey",
				FkeyResolver: dir.Arg("fkey_resolver"),
				FkeyField:    dir.Arg("fkey_name"),
			})
		}

		prefix := filepath.Join("Admin", o.NameTitle(), "Asset")

		{
			r := "Object config"
			rendered, err := gen.ExecuteTemplate("plugins/react_crud/templates/config.tsx.tmpl", base)
			if err != nil {
				return fmt.Errorf("failed rendering %s %s: %w", r, o.Name, err)
			}

			_, err = output.AppendOrCreate(gen.TS_FRONTEND_SKEL, prefix+"Config", rendered)
			if err != nil {
				return fmt.Errorf("failed appending %s %s: %w", r, o.Name, err)
			}
		}

		{
			r := "Object grid TS"
			rendered, err := gen.ExecuteTemplate("plugins/react_crud/templates/grid.tsx.tmpl", base)
			if err != nil {
				return fmt.Errorf("failed rendering %s %s: %w", r, o.Name, err)
			}

			_, err = output.AppendOrCreate(gen.TS_FRONTEND_GEN, prefix+"Grid", rendered)
			if err != nil {
				return fmt.Errorf("failed appending %s %s: %w", r, o.Name, err)
			}
		}

		{
			r := "Object grid SASS"
			rendered, err := gen.ExecuteTemplate("templates/blank.tmpl", nil)
			if err != nil {
				return fmt.Errorf("failed rendering %s %s: %w", r, o.Name, err)
			}

			_, err = output.AppendOrCreate(gen.SASS_FRONTEND_SKEL, prefix+"Grid", rendered)
			if err != nil {
				return fmt.Errorf("failed appending %s %s: %w", r, o.Name, err)
			}
		}

		{
			r := "Object crud TS"
			rendered, err := gen.ExecuteTemplate("plugins/react_crud/templates/crud.tsx.tmpl", base)
			if err != nil {
				return fmt.Errorf("failed rendering %s %s: %w", r, o.Name, err)
			}

			_, err = output.AppendOrCreate(gen.TS_FRONTEND_GEN, prefix+"CRUD", rendered)
			if err != nil {
				return fmt.Errorf("failed appending %s %s: %w", r, o.Name, err)
			}
		}

		{
			r := "Object form TS"
			rendered, err := gen.ExecuteTemplate("plugins/react_crud/templates/form.tsx.tmpl", base)
			if err != nil {
				return fmt.Errorf("failed rendering %s %s: %w", r, o.Name, err)
			}

			_, err = output.AppendOrCreate(gen.TS_FRONTEND_GEN, prefix+"Form", rendered)
			if err != nil {
				return fmt.Errorf("failed appending %s %s: %w", r, o.Name, err)
			}
		}

	}

	return nil
}
