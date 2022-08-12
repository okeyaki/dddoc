package components

import (
	"crypto/sha1"
	"fmt"
	"go/ast"

	"golang.org/x/tools/go/packages"
)

type Component interface {
	GetPkg() *packages.Package

	GetASTTypeSpec() *ast.TypeSpec
}

func GetComponentID(c Component) string {
	hash := sha1.New()
	hash.Write([]byte(GetComponentFullName(c)))

	return fmt.Sprintf("%x", hash.Sum(nil))
}

func GetComponentName(c Component) string {
	return c.GetASTTypeSpec().Name.Name
}

func GetComponentFullName(c Component) string {
	return c.GetPkg().TypesInfo.TypeOf(c.GetASTTypeSpec().Name).String()
}
