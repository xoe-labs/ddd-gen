// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package directive

import (
	. "github.com/dave/jennifer/jen"
	"go/types"
	"log"
)

func GenUnmarshalFromRepository(f *File, typeName string, publicFlds, privateFlds []*types.Var) {

	sF := shortForm(typeName)
	allFlds := append(publicFlds, privateFlds...)

	log.Printf("%s: generating 'UnmarshalFromRepository' unmarshaler\n", typeName)

	f.Commentf("UnmarshalFromRepository unmarshals %s from the repository so that non-constructable", typeName)
	f.Comment("private fields can still be initialized from (private) repository state")
	f.Comment("")
	f.Comment("Important: DO NEVER USE THIS METHOD EXCEPT FROM THE REPOSITORY")
	f.Comment("Reason: This method initializes private state, so you could corrupt the domain.")

	f.Func().Id("UnmarshalFromRepository").ParamsFunc(func(g *Group) {
		for _, field := range allFlds {
			qualType := getQualifiedJenType(field.Type(), field.Pkg())
			g.Id(field.Name()).Add(qualType)
		}
	}).Op("*").Id(typeName).BlockFunc(func(g *Group) {
		g.Id(sF).Op(":=").Id("MustNew").CallFunc(func(g *Group) {
			for _, field := range publicFlds {
				g.Id(field.Name())
			}
		})
		for _, field := range privateFlds {
			g.Id(sF).Dot(field.Name()).Op("=").Id(field.Name())
		}
		g.Return(Id(sF))
	})
}
