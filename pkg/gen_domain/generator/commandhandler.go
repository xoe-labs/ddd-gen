// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

import (
	. "github.com/dave/jennifer/jen"
	"log"
)

func GenFuncHandleStub(g *Group, typ string, entity QualId) (ident string) {
	log.Printf("%s: generating '%s()'\n", typ, Handle)

	g.Commentf("%s handles %s in the domain", Handle, typ)
	g.Comment("returns true for success or false for failure.")
	g.Commentf("record errors with %s(err).", Raise)
	g.Comment("implements application layer's CommandHandler interface.")

	g.Func().Params(
		Id(cmdShortForm(typ)).Op("*").Id(typ),
	).Id(
		Handle,
	).Params(
		Id(shortForm(entity.Id)).Op("*").Qual(entity.Qual, entity.Id),
	).Params(
		Id("bool"),
	).Block(
		Return().Id("true"),
	)
	return Handle
}

func GenFuncFacts(f *File, typ string) (ident string) {
	log.Printf("%s: generating '%s()'\n", typ, Facts)

	f.Commentf("%s returns collected domain facts", Facts)
	f.Comment("implements application layer's FactKeeper interface.")
	f.Func().Params(
		Id(cmdShortForm(typ)).Op("*").Id(typ),
	).Id(
		Facts,
	).Params().Params(
		Index().Interface(),
	).Block(
		Return().Id(cmdShortForm(typ)).Dot(
			FactsField,
		),
	)
	return Facts
}
func GenFuncErrors(f *File, typ string) (ident string) {
	log.Printf("%s: generating '%s()'\n", typ, Errors)

	f.Commentf("%s returns collected domain errors", Errors)
	f.Comment("implements application layer's ErrorKeeper interface.")
	f.Func().Params(
		Id(cmdShortForm(typ)).Op("*").Id(typ),
	).Id(
		Errors,
	).Params().Params(
		Index().Id("error"),
	).Block(
		Return().Id(cmdShortForm(typ)).Dot(
			ErrorsField,
		),
	)
	return Errors
}

func GenFuncrecordOn(f *File, typ string, entity QualId) (ident string) {
	log.Printf("%s: generating '%s()'\n", typ, RecordOn)

	f.Commentf("%s records facts and applies them to the domain", RecordOn)
	f.Func().Params(
		Id(cmdShortForm(typ)).Op("*").Id(typ),
	).Id(
		RecordOn,
	).Params(
		Id(shortForm(entity.Id)).Op("*").Qual(entity.Qual, entity.Id),
		Id("fact").Interface(),
	).Block(
		Id(shortForm(entity.Id)).Dot(
			Apply,
		).Call(
			Id("fact"),
		),
		Id(cmdShortForm(typ)).Dot(
			FactsField,
		).Op("=").Id(
			"append",
		).Call(
			Id(cmdShortForm(typ)).Dot(
				FactsField,
			),
			Id("fact"),
		),
	)
	return RecordOn
}

func GenFuncraise(f *File, typ string) (ident string) {
	log.Printf("%s: generating '%s()'\n", typ, Raise)

	f.Commentf("%s records domain errors", Raise)
	f.Func().Params(
		Id(cmdShortForm(typ)).Op("*").Id(typ),
	).Id(
		Raise,
	).Params(
		Id("err").Id("error"),
	).Block(
		Id(cmdShortForm(typ)).Dot(
			ErrorsField,
		).Op("=").Id(
			"append",
		).Call(
			Id(cmdShortForm(typ)).Dot(
				ErrorsField,
			),
			Id("err"),
		),
	)
	return Raise
}
