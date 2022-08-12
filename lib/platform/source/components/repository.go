package components

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type Repository struct {
	Pkg         *packages.Package
	ASTTypeSpec *ast.TypeSpec

	Object string
}

func (r *Repository) GetPkg() *packages.Package {
	return r.Pkg
}

func (r *Repository) GetASTTypeSpec() *ast.TypeSpec {
	return r.ASTTypeSpec
}
