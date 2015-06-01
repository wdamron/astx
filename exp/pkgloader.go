package exp

import (
	"go/ast"
	"go/parser"

	"golang.org/x/tools/go/loader"
	"golang.org/x/tools/go/types"
)

type Program struct {
	*loader.Program
	pkgNameMap map[string]*loader.PackageInfo
}

func Load(importPath string) (Program, error) {
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
	return Program{p, pkgNameMap}, err
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
