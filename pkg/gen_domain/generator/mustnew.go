// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

import (
	. "github.com/dave/jennifer/jen"
	"log"
)

func GenMustNew(f *File, typ string, publicFlds []QualField) {

	log.Printf("%s: generating '%s()'\n", typ, MustNew)

	f.Commentf("%s returns a guaranteed-to-be-valid %s or panics", MustNew, typ)

	f.Func().Id(
		MustNew,
	).ParamsFunc(func(g *Group) {
		for _, field := range publicFlds {
			g.Id(
				field.Id,
			).Add(
				field.QualTyp,
			)
		}
	}).Call(
		Op("*").Id(typ),
	).Block(
		List(
			Id(shortForm(typ)),
			Id("err"),
		).Op(":=").Id(
			Neww,
		).CallFunc(func(g *Group) {
			for _, field := range publicFlds {
				g.Id(field.Id)
			}
		}),
		If(
			Id("err").Op("!=").Id("nil"),
		).Block(
			Panic(Id("err")),
		),
		Return(
			Id(shortForm(typ)),
		),
	)
}
