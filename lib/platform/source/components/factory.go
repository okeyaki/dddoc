package components

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type Factory struct {
	Pkg         *packages.Package
	ASTTypeSpec *ast.TypeSpec

	Object string
}

func (f *Factory) GetPkg() *packages.Package {
	return f.Pkg
}

func (f *Factory) GetASTTypeSpec() *ast.TypeSpec {
	return f.ASTTypeSpec
}
