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
				fieldType := formatTypeExpr(field.Type)
				if fieldType == "" {
					continue
				}
				parsedField := StructField{
					Name: field.Names[0].Name,
					Type: fieldType,
				}
				if field.Tag != nil {
					raw := field.Tag.Value
					parsedField.RawTag = raw
					if len(raw) >= 2 {
						// Strip leading/trailing back-ticks:
						parsedField.Tag = reflect.StructTag(raw[1 : len(raw)-1])
					}
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

func formatTypeExpr(expr ast.Expr) string {
	switch V := expr.(type) {
	case *ast.Ident:
		return V.Name
	case *ast.StarExpr:
		return "*" + formatTypeExpr(V.X)
	case *ast.ArrayType:
		sz := ""
		if V.Len != nil {
			switch L := V.Len.(type) {
			case *ast.BasicLit:
				sz = L.Value
			case *ast.Ident:
				sz = L.Name
			}
		}
		return "[" + sz + "]" + formatTypeExpr(V.Elt)
	case *ast.MapType:
		return "map[" + formatTypeExpr(V.Key) + "]" + formatTypeExpr(V.Value)
	case *ast.SelectorExpr:
		return formatTypeExpr(V.X) + "." + formatTypeExpr(V.Sel)
	}
	return ""
}
