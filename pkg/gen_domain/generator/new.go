// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

import (
	. "github.com/dave/jennifer/jen"
	"log"
)

type Validation struct {
	Field  QualField
	ErrMsg string
}

func GenNew(f *File, typ string, publicFlds []QualField, validations []Validation, validatorMethod string) {

	log.Printf("%s: generating '%s()'\n", typ, Neww)

	f.Commentf("%s returns a guaranteed-to-be-valid %s or an error", Neww, typ)

	f.Func().Id(
		Neww,
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
		Id("error"),
	).BlockFunc(func(g *Group) {
		for _, fld := range validations {
			field := fld.Field
			errMsg := fld.ErrMsg
			// ... build "if <field> == nil { return nil, errors.New("<errMsg>") }"
			g.If(
				Qual(
					"reflect",
					"ValueOf",
				).Call(
					Id(field.Id),
				).Dot("IsZero").Call(),
			).Block(
				Return(
					Id("nil"),
					Qual(
						"errors",
						"New",
					).Call(
						Lit(errMsg),
					),
				),
			)
		}
		g.Id(shortForm(typ)).Op(":=").Op("&").Id(
			typ,
		).Values(
			DictFunc(func(d Dict) {
				for _, field := range publicFlds {
					d[Id(field.Id)] = Id(field.Id)
				}
			}),
		)
		if validatorMethod != "" {
			g.If(
				Id("err").Op(":=").Id(
					shortForm(typ),
				).Dot(
					validatorMethod,
				).Call(),
				Id("err").Op("!=").Id("nil"),
			).Block(
				Return(
					Id("nil"),
					Id("err"),
				),
			)
		}
		g.Return(
			Id(shortForm(typ)),
			Id("nil"),
		)
	})
}
