// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

import (
	. "github.com/dave/jennifer/jen"
	"log"
)

type EqualFld struct {
	Field       QualField
	IsDeepEqual bool
}

func GenEqual(f *File, typ string, flds []EqualFld) {

	log.Printf("%s: generating '%s()'\n", typ, Equal)

	f.Commentf("%s answers whether v is equivalent to %s", Equal, shortForm(typ))
	f.Commentf("Always returns false if v is not a %s", typ)

	f.Func().Params(
		Id(shortForm(typ)).Id(typ),
	).Id(
		Equal,
	).Params(
		Id("v").Interface(),
	).Bool().BlockFunc(func(g *Group) {
		g.List(
			Id("other"),
			Id("ok"),
		).Op(":=").Id("v").Assert(Id(typ))
		g.If(
			Op("!").Id("ok"),
		).Block(
			Return(
				Id("false"),
			),
		)

		for _, fld := range flds {
			field := fld.Field
			if fld.IsDeepEqual {
				g.If(
					Op("!").Qual(
						"reflect",
						"DeepEqual",
					).Call(
						Id(shortForm(typ)).Dot(field.Id),
						Id("other").Dot(field.Id),
					).Block(
						Return(
							Id("false"),
						),
					),
				)
			} else {
				g.If(
					Id(shortForm(typ)).Dot(field.Id).Op("!=").Id("other").Dot(field.Id),
				).Block(
					Return(
						Id("false"),
					),
				)
			}
		}
		g.Return(
			Id(shortForm(typ)),
		)
	})
}
