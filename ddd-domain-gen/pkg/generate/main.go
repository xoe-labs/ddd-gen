// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generate

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	. "github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

func Main(sourceTypeName string) error {

	// Get the package of the file with go:generate comment
	goPackage := os.Getenv("GOPACKAGE")
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	// Build the target file name
	goFile := os.Getenv("GOFILE")
	ext := filepath.Ext(goFile)
	baseFilename := goFile[0 : len(goFile)-len(ext)]
	targetFilename := baseFilename + "_gen.go"

	// Remove existing target file (before loading the package)
	if _, err := os.Stat(targetFilename); err == nil {
		if err := os.Remove(targetFilename); err != nil {
			return err
		}
	}

	// Inspect package and use type checker to infer imported types
	pkg, err := loadPackage(path)
	if err != nil {
		return err
	}

	// Lookup the given source type name in the package declarations
	obj := pkg.Types.Scope().Lookup(sourceTypeName)
	if obj == nil {
		return fmt.Errorf("%s not found in declared types of %s",
			sourceTypeName, pkg)
	}

	// We check if it is a declared type
	if _, ok := obj.(*types.TypeName); !ok {
		return fmt.Errorf("%v is not a named type", obj)
	}
	// We expect the underlying type to be a struct
	structType, ok := obj.Type().Underlying().(*types.Struct)
	if !ok {
		return fmt.Errorf("type %v is not a struct", obj)
	}

	// Generate code using jennifer
	err = generate(goPackage, targetFilename, sourceTypeName, structType)
	if err != nil {
		return err
	}
	return nil
}

// StructTag Key
var (
	structTagGenKey = "gen"
	structTagDDDKey = "ddd"
)

// A simple regexp pattern to match tag values
var (
	structRequiredTagPattern  = regexp.MustCompile(`required'([^']+)'`)
	structPrivateTagPattern   = regexp.MustCompile(`private`)
	structGenGetterTagPattern = regexp.MustCompile(`getter`)
)

func generate(goPackage, targetFilename, sourceTypeName string, structType *types.Struct) error {

	// Start a new file in this package
	// return fmt.Errorf(goPackage)
	f := NewFile(goPackage)

	// Add a package comment, so IDEs detect files as generated
	f.PackageComment("Code generated by ddd-domain-gen, DO NOT EDIT.")

	f.Comment("Generators ...")
	f.Line()

	// 1. define code region variables
	var (
		publicFields      []*types.Var
		privateFields     []*types.Var
		genGetterFields   []*types.Var
		publicParams      []Code
		privateParams     []Code
		allParams         []Code
		publicValidations []Code
	)

	// 2. iterate over struct fields and populate those variables
	for i := 0; i < structType.NumFields(); i++ {
		field := structType.Field(i)
		tag := reflect.StructTag(structType.Tag(i))

		// 2.1 match default getter creation to fields
		if structTagGenKeyValue, ok := tag.Lookup(structTagGenKey); ok {
			if matches := structGenGetterTagPattern.FindStringSubmatch(structTagGenKeyValue); matches != nil {
				genGetterFields = append(genGetterFields, field)
			}
		}

		// 2.2 separate private and public fields
		var private bool
		var requiredMatches []string
		if structTagDDDKeyValue, ok := tag.Lookup(structTagDDDKey); ok {
			if matches := structPrivateTagPattern.FindStringSubmatch(structTagDDDKeyValue); matches != nil {
				private = true
			}
			requiredMatches = structRequiredTagPattern.FindStringSubmatch(structTagDDDKeyValue)
		}
		if private {
			privateParams = append(privateParams, Id(field.Name()).Add(
				getQualifiedType(field.Type().String()),
			))
			privateFields = append(privateFields, field)
		} else {
			publicParams = append(publicParams, Id(field.Name()).Add(
				getQualifiedType(field.Type().String()),
			))
			publicFields = append(publicFields, field)
		}

		// 2.2 generate required validation code (error if also private)
		if requiredMatches != nil {
			if private {
				return fmt.Errorf("private field %s cannot be required", field.Name())
			}
			errMsg := requiredMatches[1]

			// ... build "if <field> == nil { return nil, errors.New("<errMsg>") }"
			publicValidations = append(publicValidations,
				If(Id(field.Name()).Op("==").Nil()).Block(
					Return(Nil(), Qual("errors", "New").Call(Lit(errMsg))),
				),
			)
		}
	}
	allParams = append(publicParams, privateParams...)

	// 3. assemble methods ...

	sF := shortForm(sourceTypeName)

	// -- Add New() constructor
	f.Commentf("New returns a guaranteed-to-be-valid %s or an error", sourceTypeName)
	f.Func().Id("New").Params(
		publicParams...,
	).Call(
		Op("*").Id(sourceTypeName),
		Error(),
	).BlockFunc(func(g *Group) {
		for _, code := range publicValidations {
			g.Add(code)
		}
		g.Return(Op("&").Id(sourceTypeName).Values(
			DictFunc(func(d Dict) {
				for _, fld := range publicFields {
					d[Id(fld.Name())] = Id(fld.Name())
				}
			}),
		), Nil())
	})

	// -- Add MustNew() constructor
	f.Commentf("MustNew returns a guaranteed-to-be-valid %s or panics", sourceTypeName)
	f.Func().Id("MustNew").Params(
		publicParams...,
	).Call(
		Op("*").Id(sourceTypeName),
	).Block(
		Id(sF).Op(",").Err().Op(":=").Id("New").CallFunc(func(g *Group) {
			for _, fld := range publicFields {
				g.Id(fld.Name())
			}
		}),
		If(Err().Op("!=").Nil()).Block(
			Panic(Err()),
		),
		Return(Id(sF)),
	)

	f.Comment("Marshalers ...")
	f.Line()

	// -- Add UnmarshalFromRepository() unmarshaler
	f.Commentf("UnmarshalFromRepository unmarshals %s from the repository so that non-constructable", sourceTypeName)
	f.Comment("private fields can still be initialized from (private) repository state")
	f.Comment("")
	f.Comment("Important: DO NEVER USE THIS METHOD EXCEPT FROM THE REPOSITORY")
	f.Comment("Reason: This method initializes private state, so you could corrupt the domain.")
	f.Func().Id("UnmarshalFromRepository").Params(
		allParams...,
	).Op("*").Id(sourceTypeName).BlockFunc(func(g *Group) {
		g.Id(sF).Op(":=").Id("MustNew").CallFunc(func(g *Group) {
			for _, fld := range publicFields {
				g.Id(fld.Name())
			}
		})
		for _, fld := range privateFields {
			g.Id(sF).Dot(fld.Name()).Op("=").Id(fld.Name())
		}
		g.Return(Id(sF))
	})

	f.Comment("Getters ...")
	f.Line()
	for _, fld := range genGetterFields {
		fN := strings.Title(fld.Name())
		f.Commentf("%s returns %s value", fN, fld.Name())
		f.Func().Params(
			Id(sF).Op("*").Id(sourceTypeName),
		).Id(fN).Params().Add(
			getQualifiedType(fld.Type().String()),
		).Block(
			Return(Id(sF).Dot(fld.Name())),
		)
	}

	// Write generated file
	return f.Save(targetFilename)
}

func loadPackage(path string) (*packages.Package, error) {
	cfg := &packages.Config{Mode: packages.NeedTypes | packages.NeedImports}
	pkgs, err := packages.Load(cfg, path)
	if err != nil {
		return nil, fmt.Errorf("loading packages for inspection: %v", err)
	}
	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	return pkgs[0], nil
}

func shortForm(typeName string) string {
	return strings.ToLower(string(typeName[0]))
}

func getQualifiedType(s string) *Statement {
	frst := ""
	last := s[strings.LastIndex(s, ".")+1:]
	if last != s {
		if string(s[0]) == "*" {
			// remove first character (* - ptr)
			frst = s[1:strings.LastIndex(s, ".")]
			return Op("*").Qual(frst, last)
		}
		frst = s[:strings.LastIndex(s, ".")]
	}
	return Qual(frst, last)

}
