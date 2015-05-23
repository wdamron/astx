package astx

import (
	"go/ast"
	"go/token"
	"reflect"
)

// Struct represents a struct type definition parsed from a Go source file.
type Struct struct {
	Name     string
	Comments []string
	Fields   []StructField
}

// StructField represents a field definition, within a struct type definition,
// parsed from a Go source file.
type StructField struct {
	Name   string
	Type   string
	Tag    reflect.StructTag
	RawTag string
}

// ParseFileStructs parses all struct definitions within the file at the given path.
func ParseFileStructs(fset *token.FileSet, f *ast.File) []Struct {
	parsedStructs := []Struct{}
	commentMap := ast.NewCommentMap(fset, f, f.Comments)

	for _, decl := range f.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		for _, spec := range genDecl.Specs {
			typeSpec, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}

			structType, ok := typeSpec.Type.(*ast.StructType)
			if !ok {
				continue
			}

			parsedStruct := Struct{Name: typeSpec.Name.Name}

			commentGroups := commentMap[genDecl]
			if commentGroups != nil {
				parsedStruct.Comments = []string{}
				for _, group := range commentGroups {
					for _, comment := range group.List {
						parsedStruct.Comments = append(parsedStruct.Comments, comment.Text)
					}
				}
			}

			if structType.Fields.List != nil {
				parsedStruct.Fields = make([]StructField, 0, len(structType.Fields.List))
			}
			for _, field := range structType.Fields.List {
				var fieldType string
				if tname, ok := field.Type.(*ast.Ident); ok {
					fieldType = tname.Name
				} else if tname, ok := field.Type.(*ast.StarExpr); ok {
					fieldType = "*" + tname.X.(*ast.Ident).Name
				}
				parsedField := StructField{
					Name: field.Names[0].Name,
					Type: fieldType,
				}
				if field.Tag != nil {
					parsedField.RawTag = field.Tag.Value
					parsedField.Tag = reflect.StructTag(field.Tag.Value)
				}
				parsedStruct.Fields = append(parsedStruct.Fields, parsedField)
			}

			parsedStructs = append(parsedStructs, parsedStruct)
		}
	}

	if len(parsedStructs) == 0 {
		return nil
	}

	return parsedStructs
}
