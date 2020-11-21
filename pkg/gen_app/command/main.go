// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package command

import (
	"fmt"
	"os"
	"path"
	"strings"
	"unicode"

	"log"

	"go/ast"
	"go/build"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"

	"github.com/xoe-labs/ddd-gen/pkg/gen_app/directive"
)

const (
	IdentifiableIdentiferMethodName    = "Identifier"
	PoliceableUserMethodName           = "User"
	PoliceableElevationTokenMethodName = "ElevationToken"
)

type Config struct {
	// Domain aggregate entity
	aggEntityStruct string
	// Application defined interfaces
	policeableInterface   string
	identifiableInterface string
	repositoryInterface   string
	policerInterface      string
	// Error constructors
	authorizationErrorNew  string
	identificationErrorNew string
	repositoryErrorNew     string
	domainErrorNew         string
}

func NewConfig(aggEntityStruct, policeableInterface, identifiableInterface, repositoryInterface, policerInterface,
	authorizationErrorNew, identificationErrorNew, repositoryErrorNew, domainErrorNew string) (*Config, error) {

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
	if authorizationErrorNew == "" {
		return nil, fmt.Errorf("authorizationErrorNew not set")
	}
	if identificationErrorNew == "" {
		return nil, fmt.Errorf("identificationErrorNew not set")
	}
	if repositoryErrorNew == "" {
		return nil, fmt.Errorf("repositoryErrorNew not set")
	}
	if domainErrorNew == "" {
		return nil, fmt.Errorf("domainErrorNew not set")
	}
	return &Config{
		aggEntityStruct:        aggEntityStruct,
		policeableInterface:    policeableInterface,
		identifiableInterface:  identifiableInterface,
		repositoryInterface:    repositoryInterface,
		policerInterface:       policerInterface,
		authorizationErrorNew:  authorizationErrorNew,
		identificationErrorNew: identificationErrorNew,
		repositoryErrorNew:     repositoryErrorNew,
		domainErrorNew:         domainErrorNew,
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
	_ = os.Mkdir(genPath, os.ModeDir|os.ModePerm)

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

	if !isValidQualId(conf.identifiableInterface) {
		return fmt.Errorf("%s is not a valid full qualifier", conf.identifiableInterface)
	}
	if !isValidQualId(conf.policerInterface) {
		return fmt.Errorf("%s is not a valid full qualifier", conf.policerInterface)
	}
	if !isValidQualId(conf.aggEntityStruct) {
		return fmt.Errorf("%s is not a valid full qualifier", conf.aggEntityStruct)
	}
	if !isValidQualId(conf.policeableInterface) {
		return fmt.Errorf("%s is not a valid full qualifier", conf.policeableInterface)
	}
	if !isValidQualId(conf.repositoryInterface) {
		return fmt.Errorf("%s is not a valid full qualifier", conf.repositoryInterface)
	}
	if !isValidQualId(conf.authorizationErrorNew) {
		return fmt.Errorf("%s is not a valid full qualifier", conf.authorizationErrorNew)
	}
	if !isValidQualId(conf.identificationErrorNew) {
		return fmt.Errorf("%s is not a valid full qualifier", conf.identificationErrorNew)
	}
	if !isValidQualId(conf.repositoryErrorNew) {
		return fmt.Errorf("%s is not a valid full qualifier", conf.repositoryErrorNew)
	}
	if !isValidQualId(conf.domainErrorNew) {
		return fmt.Errorf("%s is not a valid full qualifier", conf.domainErrorNew)
	}

	// Inspect Interfaces, look for given methods and detect their qualified return type
	identifiableInterface := splitQual(conf.identifiableInterface)
	policeableInterface := splitQual(conf.policeableInterface)

	identifierTyp, err := inspectInterfaceMethodForRetTyp(identifiableInterface, IdentifiableIdentiferMethodName)
	if err != nil {
		return fmt.Errorf("unable to infer identifier type: %w", err)
	}
	userTyp, err := inspectInterfaceMethodForRetTyp(policeableInterface, PoliceableUserMethodName)
	if err != nil {
		return fmt.Errorf("unable to infer user type: %w", err)
	}
	elevationTokenTyp, err := inspectInterfaceMethodForRetTyp(policeableInterface, PoliceableElevationTokenMethodName)
	if err != nil {
		return fmt.Errorf("unable to infer elevation token type: %w", err)
	}

	// Generate code using jennifer
	err = generate(genPath, sourceTypeName, structType, directive.ParsedConfig{
		AggEntityStruct:     splitQual(conf.aggEntityStruct),
		RepositoryInterface: splitQual(conf.repositoryInterface),
		PolicerInterface:    splitQual(conf.policerInterface),

		IdentifiableInterface: identifiableInterface,
		PoliceableInterface:   policeableInterface,

		IdentifierTyp:     identifierTyp,
		UserTyp:           userTyp,
		ElevationTokenTyp: elevationTokenTyp,

		AuthorizationErrorNew:  splitQual(conf.authorizationErrorNew),
		IdentificationErrorNew: splitQual(conf.identificationErrorNew),
		RepositoryErrorNew:     splitQual(conf.repositoryErrorNew),
		DomainErrorNew:         splitQual(conf.domainErrorNew),
	})
	if err != nil {
		return err
	}
	return nil
}

func newTrue() *bool {
	b := true
	return &b
}

func newFalse() *bool {
	b := false
	return &b
}

func inspectInterfaceMethodForRetTyp(fq directive.QualId, funId string) (qualT directive.QualId, err error) {
	fset := token.NewFileSet()
	ctx := build.Default
	pkg, err := ctx.Import(fq.Qual, ".", 0)
	if err != nil {
		return directive.QualId{}, err
	}

	log.Printf("... inspecting:\t%s\n", pkg.Dir)
	// log.Println("pkgName:", pkg.Name)

	astPkgs, err := parser.ParseDir(fset, pkg.Dir, nil, 0)
	if err != nil {
		return directive.QualId{}, err
	}
	astPkg := astPkgs[pkg.Name]

	var (
		breakGlass  *bool
		pkgs        []*packages.Package
		importMap   map[string]string = map[string]string{}
		glassBroken func() bool       = func() bool { return !*breakGlass }
	)

	breakGlass = newFalse()

	// Map package names or aliases to their import paths
	ast.Inspect(astPkg, func(n ast.Node) bool {
		var (
			ident *ast.Ident
			lit   *ast.BasicLit
		)
		switch nt := n.(type) {
		case *ast.ImportSpec:
			// got an import
			ident = nt.Name
			lit = nt.Path
			pkgPath := strings.Trim(lit.Value, `"`)
			if ident == nil {
				// log.Printf("\t... importing for analysis: %v\n", pkgPath)
				// got no alias -> need to import to reliably determin package name
				pkgs, err = packages.Load(&packages.Config{Mode: packages.NeedName}, pkgPath)
				importMap[pkgs[0].Name] = pkgPath
				log.Printf("\t... analyzed package import: %s\n", pkgs[0].Name)
			} else {
				importMap[ident.Name] = pkgPath
				log.Printf("\t... analyzed package import: %s\n", ident.Name)
			}

		}
		return glassBroken()
	})

	if err != nil {
		log.Printf("returning error: %v\n", err)
		return directive.QualId{}, err
	}

	if !ast.PackageExports(astPkg) {
		return directive.QualId{}, fmt.Errorf("package %s has no exported nodes", pkg.Dir)
	}

	// reinitialize
	qualT = directive.QualId{}
	breakGlass = newFalse()

	// find qualT retru type for our method
	ast.Inspect(astPkg, func(n ast.Node) bool {
		var expr ast.Expr
		switch nt := n.(type) {
		case *ast.TypeSpec:
			// got a type
			if nt.Name.String() != fq.Id {
				// it not "our" type
				return false
			}
			expr = nt.Type
			it, ok := expr.(*ast.InterfaceType)
			if !ok {
				// was not an interface type
				err = fmt.Errorf("%s is not an interface", nt.Name.String())
				breakGlass = newTrue()
				return false
			}

			log.Println("\t... found Iface:", nt.Name.String())

			// was "our" interface
			for _, m := range it.Methods.List {
				if m.Names[0].Name == funId {
					// had the method we look for
					log.Println("\t... found IfaceMethod:", m.Names[0].String())
					expr = m.Type
					ft := expr.(*ast.FuncType)
					fr := ft.Results
					if len(fr.List) != 1 {
						// had more than one result
						err = fmt.Errorf("%s has not exactly one return type: %v", funId, fr.List)
						breakGlass = newTrue()
						return false
					}
					frf := fr.List[0] // our only return type "field"
					ast.Inspect(frf, func(n ast.Node) bool {
						var expr ast.Expr
						switch nt := n.(type) {
						case *ast.Ident:
							// got a single identifier
							qualT.Id = nt.Name
							qualT.Qual = ""
							log.Printf("\t... found qualified type Id: '%s' Qual: '%s'\n", qualT.Id, qualT.Qual)
							// found what we were looking for
							breakGlass = newTrue()
							return false
						case *ast.SelectorExpr:
							// log.Println("\t... found SelectorExpr:", nt)
							// got a selector expr (= "package.Name" / "X.Sel")
							qualT.Id = nt.Sel.Name
							expr = nt.X

							id, ok := expr.(*ast.Ident)
							if !ok {
								// part before the last dot was not an identifier, we don't know what to do
								err = fmt.Errorf("%s return type is not of format package.Name - not supported: %v", funId, nt)
								breakGlass = newTrue()
								return false
							}
							iPath := importMap[id.Name]
							if iPath == "" {
								// package name/alias was not in import paths
								err = fmt.Errorf("%s return package name/alias was not found in imports: %v", funId, id)
								breakGlass = newTrue()
								return false
							}
							qualT.Qual = iPath
							// found what we were looking for
							log.Printf("\t... found qualified type Id: '%s' Qual: '%s'\n", qualT.Id, qualT.Qual)
							breakGlass = newTrue()
							return false
						}
						return glassBroken()
					})
				}
			}
		}
		return glassBroken()
	})

	if err != nil {
		return directive.QualId{}, err
	}
	// log.Println("qualT:", qualT)
	return qualT, nil
}

func isValidQualId(s string) bool {
	idx := strings.LastIndex(s, ".")
	if idx != -1 {
		id := s[strings.LastIndex(s, ".")+1:] // suggested identifier
		if strings.Index(id, "/") == -1 {     // no '/' in suggested identifier
			return unicode.IsUpper(rune(id[0])) // starts with upper case (is exported)
		}
	}
	return false
}

func splitQual(s string) directive.QualId {
	imp := s[:strings.LastIndex(s, ".")]
	id := s[strings.LastIndex(s, ".")+1:]
	return directive.QualId{
		Qual: imp,
		Id:   id,
	}
}
