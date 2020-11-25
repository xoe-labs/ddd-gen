// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package gen_app

import (
	"fmt"
	"os"
	"path"
	"strings"

	"go/ast"
	"go/parser"
	"go/token"
	"go/types"
)

func Gen(sourceTypeName string, useFactStorage bool, conf *Config) error {

	// Get the package of the file with go:generate comment
	goPackage := os.Getenv("GOPACKAGE")
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	ifacesPath := path.Join(cwd, "../")

	// Generate interfaces using jennifer
	err = generateIfaces(ifacesPath, useFactStorage, &conf.Objects, &conf.Adapters)
	if err != nil {
		return err
	}

	// Generate docfile before loading package
	docFile := path.Join(cwd, "doc.go")
	generateDoc(docFile)

	// // Build the target file name for generated code
	invokingFile := path.Join(cwd, os.Getenv("GOFILE"))

	// Inspect package and use type checker to infer imported types
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, invokingFile, nil, 0)
	if err != nil {
		return err
	}
	astObj := astFile.Scope.Objects[sourceTypeName]
	if astObj.Kind != ast.Typ {
		return fmt.Errorf("%s is not a type declaration", sourceTypeName)
	}
	astTypeSpec := astObj.Decl.(*ast.TypeSpec)
	astStructType := astTypeSpec.Type.(*ast.StructType)
	var fields []*types.Var
	var tags []string

	for _, field := range astStructType.Fields.List {
		fields = append(fields, types.NewField(
			field.Pos(),
			types.NewPackage(cwd, goPackage),
			field.Names[0].Name,
			types.Default(nil),
			false,
		))
		tags = append(tags, strings.Trim(field.Tag.Value, "`"))
	}

	structType := types.NewStruct(fields, tags)

	// Generate code using jennifer
	err = analyzeStructAndGenerateCommandWrappers(cwd, sourceTypeName, useFactStorage, structType, conf.Adapters, conf.Objects, conf.Errors)
	if err != nil {
		return err
	}
	return nil
}

