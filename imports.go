package astx

import (
	"go/ast"
)

type Import struct {
	Name     string
	Path     string
	Doc      []string
	Comments []string
}

func ParseFileImports(f *ast.File) []Import {
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
		if spec.Doc.List != nil {
			imp.Doc = make([]string, 0, len(spec.Doc.List))
		}
		for _, doc := range spec.Doc.List {
			imp.Doc = append(imp.Doc, doc.Text)
		}
		if spec.Comment.List != nil {
			imp.Comments = make([]string, 0, len(spec.Comment.List))
		}
		for _, comment := range spec.Comment.List {
			imp.Comments = append(imp.Comments, comment.Text)
		}
		imports = append(imports, imp)
	}
	return imports
}
