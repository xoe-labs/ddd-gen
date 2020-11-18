// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package directive

import (
	. "github.com/dave/jennifer/jen"
	"go/types"
	"log"
)

type Validation struct {
	Field  *types.Var
	ErrMsg string
}

func GenNew(f *File, typeName string, publicFlds []*types.Var, validations []Validation, validatorMethod string) {

	sF := shortForm(typeName)

	log.Printf("%s: generating 'New' constructor\n", typeName)

	f.Commentf("New returns a guaranteed-to-be-valid %s or an error", typeName)

	f.Func().Id("New").ParamsFunc(func(g *Group) {
		for _, field := range publicFlds {
			qualType := getQualifiedJenType(field.Type(), field.Pkg())
			g.Id(field.Name()).Add(qualType)
		}
	}).Call(
		Op("*").Id(typeName),
		Error(),
	).BlockFunc(func(g *Group) {
		for _, fld := range validations {
			field := fld.Field
			errMsg := fld.ErrMsg
			// ... build "if <field> == nil { return nil, errors.New("<errMsg>") }"
			g.If(
				Qual("reflect", "ValueOf").Call(
					Id(field.Name()),
				).Dot("IsZero").Call(),
			).Block(
				Return(Nil(), Qual("errors", "New").Call(Lit(errMsg))),
			)
		}
		g.Id(sF).Op(":=").Op("&").Id(typeName).Values(
			DictFunc(func(d Dict) {
				for _, field := range publicFlds {
					d[Id(field.Name())] = Id(field.Name())
				}
			}),
		)
		if validatorMethod != "" {
			g.If(Err().Op(":=").Id(sF).Dot(validatorMethod).Call()).Op(";").Err().Op("!=").Nil().Block(
				Return(Nil(), Err()),
			)
		}
		g.Return(Id(sF), Nil())
	})
}
