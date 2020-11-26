// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package gen_domain

import (
	"bytes"
	"fmt"
	"log"
	"os"

	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"
	"path/filepath"

	"github.com/dave/jennifer/jen"

	"golang.org/x/tools/go/packages"
)

var (
	goFile         string
	goPackagePath  string
	goPackage      string
	baseFilename   string
	targetFilename string
	cwd            string
	pkg            *packages.Package
)

func GenEntity(typ, validatorMethod string) (err error) {
	var (
		f          *jen.File
		ok         bool
		obj        types.Object
		typStruct  *types.Struct
		sourceFile *os.File
	)

	err = initMain()
	if err != nil {
		return err
	}

	// Lookup the given source type name in the package declarations
	obj = pkg.Types.Scope().Lookup(typ)
	if obj == nil {
		return fmt.Errorf("%s not found in declared types of %s",
			typ, pkg)
	}

	// We check if it is a declared type
	if _, ok = obj.(*types.TypeName); !ok {
		return fmt.Errorf("%v is not a named type", obj)
	}
	// We expect the underlying type to be a struct
	typStruct, ok = obj.Type().Underlying().(*types.Struct)
	if !ok {
		return fmt.Errorf("type %v is not a struct", obj)
	}

	log.Printf("Generating code for: %s.%s\n", goPackagePath, typ)
	f = jen.NewFilePathName(goPackagePath, goPackage)
	// Generate code using jennifer
	err = generateEntityHelperMethods(f, typ, validatorMethod, typStruct)
	if err != nil {
		return err
	}
	err = f.Save(targetFilename)
	if err != nil {
		return err
	}
	found, err := inspectPackageForMethod(goFile, typ, "Apply")
	if err != nil {
		return err
	}
	if found {
		return nil
	}

	// Add Apply method stub on the invoking file
	s := jen.CustomFunc(jen.Options{Multi: true}, func(g *jen.Group) {
		generateEntityApplyStub(g, typ)
	})
	buf := &bytes.Buffer{}
	err = s.Render(buf)
	if err != nil {
		return err
	}
	sourceFile, err = os.OpenFile(goFile,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer sourceFile.Close()
	if _, err := sourceFile.WriteString(buf.String()); err != nil {
		return err
	}
	return nil
}

func GenCommandHandler(cfg *Config) (err error) {
	var (
		f *jen.File
	)
	err = initMain()
	if err != nil {
		return err
	}
	log.Printf("Generating code for: %s.%s\n", goPackagePath, cfg.Typ)
	f = jen.NewFilePathName(goPackagePath, goPackage)
	// Generate code using jennifer
	generateCommandHelperMethods(f, cfg.Typ, cfg.Entity)
	err = f.Save(targetFilename)
	if err != nil {
		return err
	}

	// Handle already in source file
	found, err := inspectPackageForMethod(goFile, cfg.Typ, "Handle")
	if err != nil {
		return err
	}
	if found {
		return nil
	}

	// Add Handle method stub on the invoking file
	f = jen.NewFilePath(goPackagePath)
	generateCommandHandleStub(f, cfg.Typ, cfg.Entity)
	return f.Save(baseFilename + "_handle.go")
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

func initMain() (err error) {
	// Get the package of the file with go:generate comment
	goPackage = os.Getenv("GOPACKAGE")
	cwd, err = os.Getwd()
	if err != nil {
		return err
	}

	// Build the target file name
	goFile = os.Getenv("GOFILE")
	ext := filepath.Ext(goFile)
	baseFilename = goFile[0 : len(goFile)-len(ext)]
	targetFilename = baseFilename + "_gen.go"

	// Remove existing target file (before loading the package)
	if _, err = os.Stat(targetFilename); err == nil {
		if err = os.Remove(targetFilename); err != nil {
			return err
		}
	}
	// Inspect package and use type checker to infer imported types
	pkg, err = loadPackage(cwd)
	if err != nil {
		return err
	}

	goPackagePath = pkg.Types.Path()
	return nil
}

func inspectPackageForMethod(file, typ, method string) (found bool, err error) {
	fset := token.NewFileSet()
	ctx := build.Default
	pkg, err := ctx.Import(file, ".", 0)
	if err != nil {
		return false, err
	}

	log.Printf("... inspecting: %s: %s for %s\n", pkg.Dir, typ, method)
	// log.Println("pkgName:", pkg.Name)

	astPkgs, err := parser.ParseDir(fset, pkg.Dir, nil, 0)
	if err != nil {
		return false, err
	}
	astPkg := astPkgs[pkg.Name]

	found = false

	// Map package names or aliases to their import paths
	ast.Inspect(astPkg, func(n ast.Node) bool {
		var (
			ident    *ast.Ident
			expr     ast.Expr
			starExpr *ast.StarExpr
			rcvs     []*ast.Field
		)

		switch nt := n.(type) {
		case *ast.FuncDecl:
			// got an import
			ident = nt.Name
			if ident.String() != method {
				return false
			}
			rcvs = nt.Recv.List
			for _, f := range rcvs {
				expr = f.Type
				starExpr = expr.(*ast.StarExpr)
				expr = starExpr.X
				ident = expr.(*ast.Ident)
				if ident.Name != typ {
					continue
				} else {
					found = true
				}
			}
			return false
		}
		return true

	})
	return found, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
