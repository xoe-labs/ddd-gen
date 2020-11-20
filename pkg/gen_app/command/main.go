// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package command

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

type Config struct {
	aggEntityStruct       string
	policeableInterface   string
	identifiableInterface string
	repositoryInterface   string
	policerInterface      string
}

func NewConfig(aggEntityStruct, policeableInterface, identifiableInterface, repositoryInterface, policerInterface string) (*Config, error) {

	if aggEntityStruct == "" {
		return nil, fmt.Errorf("aggEntityStruct not set")
	}
	if policeableInterface == "" {
		return nil, fmt.Errorf("policeableInterface not set")
	}
	if identifiableInterface == "" {
		return nil, fmt.Errorf("identifiableInterface not set")
	}
	if repositoryInterface == "" {
		return nil, fmt.Errorf("repositoryInterface not set")
	}
	if policerInterface == "" {
		return nil, fmt.Errorf("policerInterface not set")
	}

	return &Config{
		aggEntityStruct:       aggEntityStruct,
		policeableInterface:   policeableInterface,
		identifiableInterface: identifiableInterface,
		repositoryInterface:   repositoryInterface,
		policerInterface:      policerInterface,
	}, nil
}

func Gen(sourceTypeName string, conf Config) error {

	// Get the package of the file with go:generate comment
	goPackage := os.Getenv("GOPACKAGE")
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Attemt to create ./command directory
	genPath := path.Join(cwd, "command")
	_ = os.Mkdir(genPath, os.ModeDir | os.ModePerm)

	// Generate docfile before loading package
	docFile := path.Join(genPath, "doc.go")
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
	err = generate(genPath, sourceTypeName, structType, conf)
	if err != nil {
		return err
	}
	return nil
}

