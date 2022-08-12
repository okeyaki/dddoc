package source

import (
	"fmt"
	"go/ast"
	"regexp"
	"strings"

	"github.com/fatih/structtag"
	"github.com/okeyaki/dddoc/lib/platform"
	"github.com/okeyaki/dddoc/lib/platform/source/components"
	"golang.org/x/tools/go/packages"
)

func Analyze() ([]components.Component, error) {
	var mode packages.LoadMode
	mode |= packages.NeedName
	mode |= packages.NeedFiles
	mode |= packages.NeedImports
	mode |= packages.NeedTypes
	mode |= packages.NeedTypesSizes
	mode |= packages.NeedSyntax
	mode |= packages.NeedTypesInfo

	pkgs, err := packages.Load(
		&packages.Config{
			Mode: mode,
		},
		platform.GetConfig().GetString("packages"),
	)
	if err != nil {
		return nil, err
	}

	cs := []components.Component{}
	for _, p := range pkgs {
		pcs, err := parsePackage(p)
		if err != nil {
			return nil, err
		}

		cs = append(cs, pcs...)
	}

	return cs, nil
}

func parsePackage(pkg *packages.Package) ([]components.Component, error) {
	astGenDecls := []*ast.GenDecl{}
	for _, f := range pkg.Syntax {
		for _, d := range f.Decls {
			if d, ok := d.(*ast.GenDecl); ok {
				astGenDecls = append(astGenDecls, d)
			}
		}
	}

	astTypeSpecs := []*ast.TypeSpec{}
	for _, d := range astGenDecls {
		for _, s := range d.Specs {
			if s, ok := s.(*ast.TypeSpec); ok {
				astTypeSpecs = append(astTypeSpecs, s)
			}
		}
	}

	cs := []components.Component{}
	for _, s := range astTypeSpecs {
		c, err := parseASTTypeSpec(pkg, s)
		if err != nil {
			return nil, err
		}
		if c != nil {
			cs = append(cs, c)
		}
	}

	return cs, nil
}

func parseASTTypeSpec(pkg *packages.Package, astTypeSpec *ast.TypeSpec) (components.Component, error) {
	if platform.GetConfig().GetString("parser.ignored.name") != "" {
		matches, err := regexp.MatchString(
			platform.GetConfig().GetString("parser.ignored.name"),
			astTypeSpec.Name.Name,
		)
		if err != nil {
			return nil, err
		}
		if matches {
			return nil, err
		}
	}

	parsers := []func(pkg *packages.Package, astTypeSpec *ast.TypeSpec) (components.Component, error){
		parseEntity,
		parseFactory,
		parseRepository,
	}
	for _, parser := range parsers {
		c, err := parser(pkg, astTypeSpec)
		if err != nil {
			return nil, err
		}
		if c != nil {
			return c, nil
		}
	}

	return nil, nil
}

func parseEntity(pkg *packages.Package, astTypeSpec *ast.TypeSpec) (components.Component, error) {
	astStructType, ok := astTypeSpec.Type.(*ast.StructType)
	if !ok {
		return nil, nil
	}

	isEntity := false
	for _, f := range astStructType.Fields.List {
		if len(f.Names) == 0 {
			fmt.Println(astTypeSpec)
		}

		isMatched, err := regexp.MatchString(
			platform.GetConfig().GetString("parser.entity.id.name"),
			f.Names[0].Name,
		)
		if err != nil {
			return nil, err
		}
		if isMatched {
			isEntity = true

			break
		}
	}
	if !isEntity {
		return nil, nil
	}

	fields := []*components.EntityField{}
	for _, f := range astStructType.Fields.List {
		field, err := parseEntityField(pkg, f)
		if err != nil {
			return nil, err
		}

		fields = append(fields, field)
	}

	entity := &components.Entity{
		Pkg:         pkg,
		ASTTypeSpec: astTypeSpec,
		Fields:      fields,
	}

	return entity, nil
}

func parseFactory(pkg *packages.Package, astTypeSpec *ast.TypeSpec) (components.Component, error) {
	pattern, err := regexp.Compile(platform.GetConfig().GetString("parser.factory.name"))
	if err != nil {
		return nil, err
	}

	matches := pattern.FindAllStringSubmatch(astTypeSpec.Name.Name, -1)
	if len(matches) == 0 {
		return nil, nil
	}

	factroy := &components.Factory{
		Pkg:         pkg,
		ASTTypeSpec: astTypeSpec,
		Object:      matches[0][1],
	}

	return factroy, nil
}

func parseRepository(pkg *packages.Package, astTypeSpec *ast.TypeSpec) (components.Component, error) {
	pattern, err := regexp.Compile(platform.GetConfig().GetString("parser.repository.name"))
	if err != nil {
		return nil, err
	}

	matches := pattern.FindAllStringSubmatch(astTypeSpec.Name.Name, -1)
	if len(matches) == 0 {
		return nil, nil
	}

	repository := &components.Repository{
		Pkg:         pkg,
		ASTTypeSpec: astTypeSpec,
		Object:      matches[0][1],
	}

	return repository, nil
}

func parseEntityField(pkg *packages.Package, astField *ast.Field) (*components.EntityField, error) {
	association, err := parseEntityFieldAssociation(astField.Tag)
	if err != nil {
		return nil, err
	}

	field := &components.EntityField{
		Pkg:         pkg,
		ASTField:    astField,
		Association: association,
	}

	return field, nil
}

func parseEntityFieldAssociation(astTag *ast.BasicLit) (*components.EntityFieldAssociation, error) {
	if astTag == nil {
		return nil, nil
	}

	tags, err := structtag.Parse(strings.Trim(astTag.Value, "`"))
	if err != nil {
		return nil, err
	}

	for _, tag := range tags.Tags() {
		if tag.Key != "dddoc" {
			continue
		}

		tagName := strings.TrimSpace(tag.Name)

		isMatched, err := regexp.MatchString(
			platform.GetConfig().GetString("parser.entity.field.tag.name"),
			tagName,
		)
		if err != nil {
			return nil, err
		}
		if !isMatched {
			continue
		}

		description := tagName
		if len(tag.Options) >= 2 {
			description = strings.TrimSpace(tag.Options[1])
		}

		association := &components.EntityFieldAssociation{
			With:        strings.TrimSpace(tag.Options[0]),
			Description: description,
		}

		return association, nil
	}

	return nil, nil
}
