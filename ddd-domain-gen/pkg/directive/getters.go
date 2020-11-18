// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package directive

import (
	. "github.com/dave/jennifer/jen"
	"go/types"
	"log"
	"strings"
)

func GenGetters(f *File, typeName string, flds []*types.Var) {

	sF := shortForm(typeName)

	for _, field := range flds {
		fn := field.Name()
		fN := strings.Title(fn)
		qualType := getQualifiedJenType(field.Type(), field.Pkg())

		log.Printf("%s: generating '%s' getter\n", typeName, fN)

		f.Commentf("%s returns %s value", fN, fn)

		f.Func().Params(
			Id(sF).Op("*").Id(typeName),
		).Id(fN).Params().Add(
			qualType,
		).Block(
			Return(Id(sF).Dot(fn)),
		)
	}
}
