// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

import (
	. "github.com/dave/jennifer/jen"
	"log"
	"strings"
)

func GenGetters(f *File, typ string, flds []QualField) {

	for _, field := range flds {

		log.Printf("%s: generating '%s' getter\n", typ, strings.Title(field.Id))

		f.Commentf("%s returns %s value", strings.Title(field.Id), field.Id)

		f.Func().Params(
			Id(shortForm(typ)).Op("*").Id(typ),
		).Id(
			strings.Title(field.Id),
		).Params().Add(
			field.QualTyp,
		).Block(
			Return(
				Id(shortForm(typ)).Dot(field.Id),
			),
		)
	}
}
