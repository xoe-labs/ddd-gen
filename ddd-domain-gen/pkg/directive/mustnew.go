// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package directive

import (
	. "github.com/dave/jennifer/jen"
	"go/types"
	"log"
)

func GenMustNew(f *File, typeName string, publicFlds []*types.Var) {

	sF := shortForm(typeName)

	log.Printf("%s: generating 'MustNew' constructor\n", typeName)

	f.Commentf("MustNew returns a guaranteed-to-be-valid %s or panics", typeName)

	f.Func().Id("MustNew").ParamsFunc(func(g *Group){
		for _, field := range publicFlds {
			qualType := getQualifiedJenType(field.Type(), field.Pkg())
			g.Id(field.Name()).Add(qualType)
		}
	}).Call(
		Op("*").Id(typeName),
	).Block(
		Id(sF).Op(",").Err().Op(":=").Id("New").CallFunc(func(g *Group) {
			for _, field := range publicFlds {
				g.Id(field.Name())
			}
		}),
		If(Err().Op("!=").Nil()).Block(
			Panic(Err()),
		),
		Return(Id(sF)),
	)
}

