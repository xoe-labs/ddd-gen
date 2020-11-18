// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package entity

import (
	"fmt"
	"go/types"
	"os"
	"path/filepath"

	"golang.org/x/tools/go/packages"
)

var (
	goFile         string
	goPackagePath  string
	goPackage      string
	targetFilename string
)

func Gen(sourceTypeName, validatorMethod string) error {

	// Get the package of the file with go:generate comment
	goPackage = os.Getenv("GOPACKAGE")
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	// Build the target file name
	goFile = os.Getenv("GOFILE")
	ext := filepath.Ext(goFile)
	baseFilename := goFile[0 : len(goFile)-len(ext)]
	targetFilename = baseFilename + "_gen.go"

	// Remove existing target file (before loading the package)
	if _, err := os.Stat(targetFilename); err == nil {
		if err := os.Remove(targetFilename); err != nil {
			return err
		}
	}

	// Inspect package and use type checker to infer imported types
	pkg, err := loadPackage(cwd)
	if err != nil {
		return err
	}

	goPackagePath = pkg.Types.Path()

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
	err = generate(sourceTypeName, validatorMethod, goPackagePath, structType)
	if err != nil {
		return err
	}
	return nil
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

