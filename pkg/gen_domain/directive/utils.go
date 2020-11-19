// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package directive

import (
	"go/types"
	"log"
	"strings"

	. "github.com/dave/jennifer/jen"
)

func shortForm(typeName string) string {
	return strings.ToLower(string(typeName[0]))
}

func getQualifiedJenType(ft types.Type, pkg *types.Package) *Statement {
	switch t := ft.(type) {
	case *types.Basic:
		return Id(t.Name())
	case *types.Array:
		return Index().Add(getQualifiedJenType(t.Elem(), pkg))
	case *types.Slice:
		return Index().Add(getQualifiedJenType(t.Elem(), pkg))
	case *types.Map:
		return Map(getQualifiedJenType(t.Key(), pkg)).Add(getQualifiedJenType(t.Elem(), pkg))
	case *types.Named:
		if pkg == t.Obj().Pkg() {
			return Id(t.Obj().Name())
		}
		return Qual(t.Obj().Pkg().Path(), t.Obj().Name())
	default:
		log.Printf("Unsupported field type: %v\n", ft)
		return Nil()
	}
}

