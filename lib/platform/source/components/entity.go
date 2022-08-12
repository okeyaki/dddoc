package components

import (
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type Entity struct {
	Pkg         *packages.Package
	ASTTypeSpec *ast.TypeSpec

	Fields []*EntityField
}

func (e *Entity) GetPkg() *packages.Package {
	return e.Pkg
}

func (e *Entity) GetASTTypeSpec() *ast.TypeSpec {
	return e.ASTTypeSpec
}

type EntityField struct {
	Pkg      *packages.Package
	ASTField *ast.Field

	Association *EntityFieldAssociation
}

type EntityFieldAssociation struct {
	With        string
	Description string
}
