// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

import (
	. "github.com/dave/jennifer/jen"
	"log"
)

func GenUnmarshalFromStore(f *File, typ string, publicFlds, privateFlds []QualField) {

	allFlds := append(publicFlds, privateFlds...)

	log.Printf("%s: generating '%s()'\n", typ, UnmarshalFromStore)

	f.Commentf("%s unmarshals %s from the repository so that non-constructable", UnmarshalFromStore, typ)
	f.Comment("private fields can still be initialized from (private) repository state")
	f.Comment("")
	f.Comment("Important: DO NEVER USE THIS METHOD EXCEPT FROM THE REPOSITORY")
	f.Comment("Reason: This method initializes private state, so you could corrupt the domain.")

	f.Func().Id(
		UnmarshalFromStore,
	).ParamsFunc(func(g *Group) {
		for _, field := range allFlds {
			g.Id(
				field.Id,
			).Add(
				field.QualTyp,
			)
		}
	}).Op("*").Id(
		typ,
	).BlockFunc(func(g *Group) {
		g.Id(shortForm(typ)).Op(":=").Id(
			MustNew,
		).CallFunc(func(g *Group) {
			for _, field := range publicFlds {
				g.Id(field.Id)
			}
		})
		for _, field := range privateFlds {
			g.Id(shortForm(typ)).Dot(field.Id).Op("=").Id(field.Id)
		}
		g.Return(
			Id(shortForm(typ)),
		)
	})
}
