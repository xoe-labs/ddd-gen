// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

import (
	. "github.com/dave/jennifer/jen"
	"log"
	"strings"
)

func GenSetters(f *File, typ string, flds []QualField) {

	for _, field := range flds {
		fN := SetterPrefix + strings.Title(field.Id)

		log.Printf("%s: generating '%s()'\n", typ, fN)

		f.Commentf("%s sets %s value", fN, field.Id)

		f.Func().Params(
			Id(
				shortForm(typ),
			).Op("*").Id(
				typ,
			),
		).Id(fN).Params(
			Id(
			   field.Id,
			).Add(
				field.QualTyp,
			),
		).Block(
			Id(
				shortForm(typ),
			).Dot(
				field.Id,
			).Op("=").Id(
				field.Id,
			),
		)
	}
}
