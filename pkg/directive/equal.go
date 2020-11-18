// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package directive

import (
	. "github.com/dave/jennifer/jen"
	"go/types"
	"log"
)

type EqualFld struct {
	Field       *types.Var
	IsDeepEqual bool
}

func GenEqual(f *File, typeName string, flds []EqualFld) {

	sF := shortForm(typeName)

	log.Printf("%s: generating 'Equal' comparator\n", typeName)

	f.Commentf("Equal answers whether v is equivalent to %s", sF)
	f.Commentf("Always returns false if v is not a %s", typeName)

	f.Func().Params(
		Id(sF).Id(typeName),
	).Id("Equal").Params(
		Id("v").Interface(),
	).Bool().BlockFunc(func(g *Group) {
		g.Id("other").Op(",").Id("ok").Op(":=").Id("v").Assert(Id(typeName))
		g.If(Op("!").Id("ok")).Block(Return(False()))

		for _, fld := range flds {
			field := fld.Field
			if fld.IsDeepEqual {
				g.If(
					Op("!").Qual("reflect", "DeepEqual").Call(
						Id(sF).Dot(field.Name()),
						Id("other").Dot(field.Name()),
					).Block(
						Return(False()),
					),
				)
			} else {
				g.If(
					Id(sF).Dot(field.Name()).Op("!=").Id("other").Dot(field.Name()),
				).Block(
					Return(False()),
				)
			}
		}
		g.Return(Id(sF))
	})
}
