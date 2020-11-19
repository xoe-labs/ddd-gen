// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package directive

import (
	. "github.com/dave/jennifer/jen"
	"go/types"
	"log"
	"strings"
)

// func (f Foo) String() string {
// 	return fmt.Sprintf("%v %v %v", f.name, f.Desc, f.size)
// }

func GenStringer(f *File, typeName string, flds []*types.Var) {

	sF := shortForm(typeName)

	log.Printf("%s: generating 'String' stringer\n", typeName)

	f.Commentf("// String implements the fmt.Stringer interface and returns the native format of %s", typeName)

	f.Func().Params(
		Id(sF).Id(typeName),
	).Id("String").Params().String().Block(
		Return(Qual("fmt", "Sprintf").CallFunc(func(g *Group) {
			g.Lit(strings.Repeat("%s ", len(flds)))
			for _, field := range flds {
				g.Id(sF).Dot(field.Name())
			}
		})),
	)
}
