package exp

import (
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/types"
)

type Program struct {
	*loader.Program
	pkgNameMap map[string]*loader.PackageInfo
	fileMap    map[string]*ast.File
}

func Load(importPath string) (*Program, error) {
	cfg := &loader.Config{
		Fset:        nil,
		ParserMode:  parser.ParseComments,
		AllowErrors: true,
	}
	cfg.Import(importPath)
	p, err := cfg.Load()
	pkgNameMap := make(map[string]*loader.PackageInfo, len(p.AllPackages))
	for pkg, pkgInfo := range p.AllPackages {
		pkgNameMap[pkg.Name()] = pkgInfo
	}
	return &Program{p, pkgNameMap, nil}, err
}

func LoadDir(path string) (*Program, error) {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	cfg := &loader.Config{
		Fset:        fset,
		ParserMode:  parser.ParseComments,
		AllowErrors: true,
	}
	fileMap := make(map[string]*ast.File)
	for _, astPkg := range pkgs {
		files := make([]*ast.File, 0, len(astPkg.Files))
		for filename, af := range astPkg.Files {
			fileMap[filename] = af
			files = append(files, af)
		}
		cfg.CreateFromFiles(astPkg.Name, files...)
	}
	p, err := cfg.Load()
	pkgNameMap := make(map[string]*loader.PackageInfo, len(p.AllPackages))
	for pkg, pkgInfo := range p.AllPackages {
		pkgNameMap[pkg.Name()] = pkgInfo
	}
	return &Program{p, pkgNameMap, fileMap}, err
}

func LoadFile(filePath string) (*Program, error) {
	fset := token.NewFileSet()
	af, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	cfg := &loader.Config{
		Fset:        fset,
		ParserMode:  parser.ParseComments,
		AllowErrors: true,
	}
	cfg.CreateFromFiles(af.Name.Name, af)
	p, err := cfg.Load()
	fileMap := map[string]*ast.File{filePath: af}
	pkgNameMap := make(map[string]*loader.PackageInfo, len(p.AllPackages))
	for pkg, pkgInfo := range p.AllPackages {
		pkgNameMap[pkg.Name()] = pkgInfo
	}
	return &Program{p, pkgNameMap, fileMap}, err
}

func (p Program) PkgByName(pkgName string) *loader.PackageInfo {
	return p.pkgNameMap[pkgName]
}

func (p Program) PkgByPath(importPath string) *loader.PackageInfo {
	return p.Package(importPath)
}

func (p Program) Type(pkgName, ident string) types.Type {
	pkg := p.PkgByName(pkgName)
	if pkg == nil {
		return nil
	}
	obj := pkg.Pkg.Scope().Lookup(ident)
	if obj == nil {
		return nil
	}
	return obj.Type()
}

func (p Program) Underlying(pkgName, ident string) types.Type {
	typ := p.Type(pkgName, ident)
	if typ == nil {
		return nil
	}
	return typ.Underlying()
}

func IsStruct(t types.Type) bool {
	_, ok := t.(*types.Struct)
	return ok
}

func IsBasic(t types.Type) bool {
	_, ok := t.(*types.Basic)
	return ok
}

func typeOf(pkg *loader.PackageInfo, expr ast.Expr) types.Type {
	return pkg.Info.TypeOf(expr)
}

func objectOf(pkg *loader.PackageInfo, ident *ast.Ident) types.Object {
	return pkg.Info.ObjectOf(ident)
}

func pkgFiles(pkg *loader.PackageInfo) []*ast.File {
	return pkg.Files
}

// func fileScope(pkg *loader.PackageInfo, filePath string) types.Scope {
// 	var p ast.File
// 	p.
// 	for _, f := range pkgFiles(pkg) {

// 	}
// }
