package astx

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"reflect"
)

const (
	// Parse imports within files
	OptParseImports = 1
	// Parse structs within files
	OptParseStructs = 2
)

// Package represents a parsed Go package
type Package struct {
	Name  string
	Files []File
}

// File represents a parsed Go source file
type File struct {
	Package string
	Path    string
	AbsPath string
	Imports []Import
	Structs []Struct
}

// Import represents an import spec parsed from a Go source file
type Import struct {
	Name     string
	Path     string
	Doc      []string
	Comments []string
}

// Struct represents a struct type definition parsed from a Go source file
type Struct struct {
	Name     string
	Comments []string
	Fields   []StructField
}

// StructField represents a field definition, within a struct type definition,
// parsed from a Go source file
type StructField struct {
	Name     string
	Type     string
	Doc      []string
	Comments []string
	Tag      reflect.StructTag
	RawTag   string
	// A field definition may contain an embedded struct type definition
	StructType *Struct
}

// ParseDirOptions parses all packages at path. All options will be set to their
// default values.
func ParseDir(path string) (map[string]Package, error) {
	return ParseDirOptions(path, OptParseImports|OptParseStructs)
}

// ParseDirOptions parses all packages within path with options.
func ParseDirOptions(path string, options int) (map[string]Package, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	result := make(map[string]Package, len(pkgs))
	for pkgName, astPkg := range pkgs {
		p := Package{Name: pkgName}
		for filename, astFile := range astPkg.Files {
			f, err := parseFileOptions(filename, fset, astFile, options)
			if err != nil {
				return nil, err
			}
			p.Files = append(p.Files, *f)
		}
		result[pkgName] = p
	}
	return result, nil
}

// ParseFileOptions parses the file at path. All options will be set to their
// default values.
func ParseFile(path string) (*File, error) {
	return ParseFileOptions(path, OptParseImports|OptParseStructs)
}

// ParseFileOptions parses the file at path with options.
func ParseFileOptions(path string, options int) (*File, error) {
	fset := token.NewFileSet()
	af, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	return parseFileOptions(path, fset, af, options)
}

func parseFileOptions(path string, fset *token.FileSet, af *ast.File, options int) (*File, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	f := &File{
		Package: af.Name.Name,
		Path:    path,
		AbsPath: abs,
	}
	if options&OptParseImports != 0 {
		f.Imports = parseFileImports(af)
	}
	if options&OptParseStructs != 0 {
		f.Structs = parseFileStructs(fset, af)
	}
	return f, nil
}

// ParseSource parses src. src may be of type a string, []byte, or io.Reader.
// All options will be set to their default values.
func ParseSource(src interface{}) (*File, error) {
	return ParseSourceOptions(src, OptParseImports|OptParseStructs)
}

// ParseSourceOptions parses src with options. src may be of type a string,
// []byte, or io.Reader.
func ParseSourceOptions(src interface{}, options int) (*File, error) {
	fset := token.NewFileSet()
	af, err := parser.ParseFile(fset, "source", src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	f := &File{
		Package: af.Name.Name,
		Path:    "source",
	}
	if options&OptParseImports != 0 {
		f.Imports = parseFileImports(af)
	}
	if options&OptParseStructs != 0 {
		f.Structs = parseFileStructs(fset, af)
	}
	return f, nil
}

func parseFileImports(f *ast.File) []Import {
	if f.Imports == nil {
		return nil
	}
	imports := make([]Import, 0, len(f.Imports))
	for _, spec := range f.Imports {
		var name, path string
		if spec.Name != nil {
			name = spec.Name.Name
		}
		if spec.Path != nil {
			path = spec.Path.Value
		}
		imp := Import{
			Name: name,
			Path: path,
		}
		if spec.Doc != nil && spec.Doc.List != nil {
			imp.Doc = make([]string, 0, len(spec.Doc.List))
			for _, doc := range spec.Doc.List {
				imp.Doc = append(imp.Doc, doc.Text)
			}
		}
		if spec.Comment != nil && spec.Comment.List != nil {
			imp.Comments = make([]string, 0, len(spec.Comment.List))
			for _, comment := range spec.Comment.List {
				imp.Comments = append(imp.Comments, comment.Text)
			}
		}
		imports = append(imports, imp)
	}
	return imports
}

func parseFileStructs(fset *token.FileSet, f *ast.File) []Struct {
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
			structName := typeSpec.Name.Name
			comments := []string{}
			commentGroups := commentMap[genDecl]
			if commentGroups != nil {
				for _, group := range commentGroups {
					for _, comment := range group.List {
						comments = append(comments, comment.Text)
					}
				}
			}
			parsedStruct := parseStruct(structType)
			parsedStruct.Name = structName
			if len(comments) != 0 {
				parsedStruct.Comments = comments
			}
			parsedStructs = append(parsedStructs, *parsedStruct)
		}
	}

	if len(parsedStructs) == 0 {
		return nil
	}

	return parsedStructs
}

func parseStruct(s *ast.StructType) *Struct {
	parsedStruct := &Struct{}
	if s.Fields.List != nil {
		parsedStruct.Fields = make([]StructField, 0, len(s.Fields.List))
	}
	for _, field := range s.Fields.List {
		parsedField := StructField{}
		for i, name := range field.Names {
			parsedField.Name += name.Name
			if i != len(field.Names)-1 {
				parsedField.Name += ", "
			}
		}
		if field.Doc != nil && field.Doc.List != nil {
			parsedField.Doc = make([]string, 0, len(field.Doc.List))
			for _, doc := range field.Doc.List {
				parsedField.Doc = append(parsedField.Doc, doc.Text)
			}
		}
		if field.Comment != nil && field.Comment.List != nil {
			parsedField.Comments = make([]string, 0, len(field.Comment.List))
			for _, comment := range field.Comment.List {
				parsedField.Comments = append(parsedField.Comments, comment.Text)
			}
		}
		if field.Tag != nil {
			raw := field.Tag.Value
			parsedField.RawTag = raw
			if len(raw) >= 2 {
				// Strip leading/trailing back-ticks:
				parsedField.Tag = reflect.StructTag(raw[1 : len(raw)-1])
			}
		}
		parsedField.Type = formatTypeExpr(field.Type)
		parsedField.StructType = parseEmbeddedStructType(field.Type)
		parsedStruct.Fields = append(parsedStruct.Fields, parsedField)
	}
	return parsedStruct
}

func formatTypeExpr(expr ast.Expr) string {
	switch V := expr.(type) {
	default:
		return "?"
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
	case *ast.StructType:
		return "struct{...}"
	}
}

func parseEmbeddedStructType(expr ast.Expr) *Struct {
	switch V := expr.(type) {
	default:
		return nil
	case *ast.StructType:
		return parseStruct(V)
	case *ast.StarExpr:
		return parseEmbeddedStructType(V.X)
	}
}
