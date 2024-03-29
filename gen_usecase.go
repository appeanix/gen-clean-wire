package main

import (
	"fmt"
	"go/types"
	"strings"
)

type UseCaseFieldMeta struct {
	Field string
	Value string
}

type UseCaseMeta struct {
	Name   string
	Fields []UseCaseFieldMeta
}

const usecaseTemplate = `
// Code generated by generator, DO NOT EDIT.
package wire

{{ if gt (len .UseCases) 0 }}
import (
	. "{{.Project}}/domain/usecase"
	handler "{{.Project}}/internal/handler"
{{- if .ORM }}
	repo "{{.Project}}/internal/repo"
{{- end}}
)
{{ end }}

{{- range .UseCases }}
func new{{.Name}}(context *handler.AppContext) {{.Name}} {
	return {{ .Name }} {
{{- range .Fields }}
		{{.Field }}: {{ .Value }},
{{- end }}
	}
}
{{- end }}
`

const (
	rulerSuffix    = "Ruler"
	ruleSuffix     = "Rule"
	useCaseSuffix  = "UseCase"
	dbAdapter      = "DBAdapter"
	queryAdapter   = "QueryAdapter"
	contextAdapter = "ContextAdapter"
	commonAdapter  = "Adapter"
)

func genUseCases() {
	// Inspect package and use type checker to infer imported types
	pkg := loadPackage("")

	// Lookup the given source type name in the package declarations
	for _, usecaseName := range pkg.Types.Scope().Names() {
		// Assert suffix
		if !strings.HasSuffix(usecaseName, useCaseSuffix) {
			continue
		}

		obj := pkg.Types.Scope().Lookup(usecaseName)

		// Assert a declare type
		if _, ok := obj.(*types.TypeName); !ok {
			continue
		}

		// Assert exportable type
		if !obj.Exported() {
			continue
		}

		// Assert underlying type to be a struct
		usecaseStruct, ok := obj.Type().Underlying().(*types.Struct)
		if !ok {
			continue
		}

		// Build UseCase instance and binding
		useCase := UseCaseMeta{
			Name: usecaseName,
		}

		for i := 0; i < usecaseStruct.NumFields(); i++ {
			fieldName := usecaseStruct.Field(i).Name()

			if strings.HasSuffix(fieldName, rulerSuffix) {
				ruleName := strings.Replace(fieldName, rulerSuffix, ruleSuffix, 1)
				useCase.Fields = append(useCase.Fields, UseCaseFieldMeta{
					Field: fieldName,
					Value: fmt.Sprintf("new%s(context.AppDB, context)", ruleName),
				})

			} else if fieldName == dbAdapter {
				useCase.Fields = append(useCase.Fields, UseCaseFieldMeta{
					Field: dbAdapter,
					Value: "context.AppDB",
				})

			} else if fieldName == queryAdapter {
				useCase.Fields = append(useCase.Fields, UseCaseFieldMeta{
					Field: queryAdapter,
					Value: "repo.QueryRepo{}",
				})
				meta.ORM = true

			} else if fieldName == contextAdapter {
				useCase.Fields = append(useCase.Fields, UseCaseFieldMeta{
					Field: contextAdapter,
					Value: "context",
				})

			} else if strings.HasSuffix(fieldName, adapterSuffix) {
				// write custom wire extension function
				useCase.Fields = append(useCase.Fields, UseCaseFieldMeta{
					Field: fieldName,
					Value: fmt.Sprintf("new%s%s(context)", usecaseName, fieldName),
				})
			}
		}
		meta.UseCases = append(meta.UseCases, useCase)
	}
}
