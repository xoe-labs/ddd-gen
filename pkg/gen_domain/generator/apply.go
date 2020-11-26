// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

import (
	. "github.com/dave/jennifer/jen"
	"log"
)

func GenApplyStub(g *Group, typ string) (ident string) {
	log.Printf("%s: generating '%s()'\n", typ, Apply)

	g.Commentf("%s applies facts to %s", Apply, typ)
	g.Comment("implements application layer's entity interface.")

	g.Func().Params(
		Id(cmdShortForm(typ)).Op("*").Id(typ),
	).Id(
		Apply,
	).Params(
		Id("fact").Interface(),
	).Params(
	).Block(
		Comment("TODO: ipmlement"),
	)
	return Apply
}

