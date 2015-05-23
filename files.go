package astx

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
)

const (
	OptParseImports = 1
	OptParseStructs = 2
)

type File struct {
	Package string
	Path    string
	AbsPath string
	Imports []Import
	Structs []Struct
}

func ParseASTFile(path string) (*token.FileSet, *ast.File, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, nil, err
	}
	return fset, f, nil
}

func ParseFile(path string) (*File, error) {
	fset, af, err := ParseASTFile(path)
	if err != nil {
		return nil, err
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}
	f := &File{
		Package: af.Name.Name,
		Path:    path,
		AbsPath: abs,
	}
	f.Imports = ParseFileImports(af)
	f.Structs = ParseFileStructs(fset, af)
	return f, nil
}

func ParseFileOptions(path string, options int) (*File, error) {
	fset, af, err := ParseASTFile(path)
	if err != nil {
		return nil, err
	}
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
		f.Imports = ParseFileImports(af)
	}
	if options&OptParseStructs != 0 {
		f.Structs = ParseFileStructs(fset, af)
	}
	return f, nil
}
