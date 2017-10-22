package mapping

import (
	"go/ast"

	"github.com/koki/shorthand/inspect"
)

// IdentityMapping returns nil if this isn't a struct.
func IdentityMapping(context *inspect.Context) *MappedStruct {
	rootName := context.TypeName()
	if root, ok := context.TypeExpr().(*ast.StructType); ok {
		fields := []*MappedField{}
		for _, field := range root.Fields.List {
			for _, fieldIdent := range field.Names {
				fieldName := fieldIdent.Name
				sourceField := &SourceField{
					Name:     fieldName,
					TypeExpr: field.Type,
					Context:  context,
				}

				// This is a lie.
				mappedField := &MappedField{
					Name:       fieldName,
					OriginPath: []*SourceField{sourceField},
				}
				fields = append(fields, mappedField)
			}
		}
	}

	return nil
}
