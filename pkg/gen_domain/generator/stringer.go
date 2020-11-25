// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

import (
	. "github.com/dave/jennifer/jen"
	"log"
	"strings"
)

// func (f Foo) String() string {
// 	return fmt.Sprintf("%v %v %v", f.name, f.Desc, f.size)
// }

func GenStringer(f *File, typ string, flds []QualField) {

	log.Printf("%s: generating '%s()'\n", typ, Stringer)

	f.Commentf("%s implements the fmt.Stringer interface and returns the native format of %s", Stringer, typ)

	f.Func().Params(
		Id(shortForm(typ)).Id(typ),
	).Id(
		Stringer,
	).Params().Id(
		"string",
	).Block(
		Return(
			Qual(
				"fmt",
				"Sprintf",
			).CallFunc(func(g *Group) {
				g.Lit(strings.Repeat("%s ", len(flds)))
				for _, field := range flds {
					g.Id(shortForm(typ)).Dot(field.Id)
				}
			}),
		),
	)
}
