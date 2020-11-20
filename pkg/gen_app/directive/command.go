// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package directive

import (
	"bytes"
	"fmt"
	. "github.com/dave/jennifer/jen"
	// "log"
	"regexp"
	"strings"
)

var cmdGenCommand string = "ddd-gen app command"

// Regeneratables ...

func addCommandNotAuthorizedErr(f *File, DoSomething string) {
	f.Commentf("ErrNotAuthorizedTo%s signals that the caller is not authorized to perform %s", DoSomething, DoSomething)
	f.Null().Const().Id(
		"ErrNotAuthorizedTo"+DoSomething,
	).Op("=").Qual(
		"github.com/xoe-labs/vicidial-go/internal/common/errors",
		"AuthorizationError",
	).Call(
		Lit("ErrNotAuthorizedTo" + DoSomething),
	)
}

func addCommandNotIdentifiableErr(f *File, DoSomething string) {
	f.Commentf("Err%sNotIdentifiable signals that the command's object was not identifiable", DoSomething)
	f.Null().Const().Id(
		"Err"+DoSomething+"NotIdentifiable",
	).Op("=").Qual(
		"github.com/xoe-labs/vicidial-go/internal/common/errors",
		"IdentificationError",
	).Call(
		Lit("Err" + DoSomething + "NotIdentifiable"),
	)
}

func addCommandHandlerType(f *File, DoSomething string, adapters []struct{ Id, Qual string }) {
	f.Commentf("%sHandler knows how to perform %s", DoSomething, DoSomething)
	f.Null().Type().Id(
		DoSomething + "Handler",
	).StructFunc(func(g *Group) {
		for _, s := range adapters {
			g.Id(s.Id).Qual(splitQual(s.Qual))
		}
	})
}
func addCommandHandlerConstructor(f *File, DoSomething string, adapters []struct{ Id, Qual string }) {
	f.Commentf("New%sHandler returns %sHandler", DoSomething, DoSomething)
	f.Func().Id(
		"New" + DoSomething + "Handler",
	).ParamsFunc(func(g *Group) {
		for _, s := range adapters {
			g.Id(s.Id).Qual(splitQual(s.Qual))
		}
	}).Params(
		Op("*").Id(DoSomething + "Handler"),
	).BlockFunc(func(g *Group) {
		for _, s := range adapters {
			g.If(
				Qual(
					"reflect", "ValueOf",
				).Call(
					Id(s.Id),
				).Dot(
					"IsZero",
				).Call(),
			).Block(
				Id("panic").Call(Lit("no '" + s.Id + "' provided!")),
			)
		}
		g.Return().Op("&").Id(
			DoSomething + "Handler",
		).ValuesFunc(func(g *Group) {
			for _, s := range adapters {
				g.Id(s.Id).Op(":").Id(s.Id)
			}
		})
	})
}
func addCommandFuncHandle(f *File, DoSomething string, withPolicy bool, aggEntity string) {
	entityImp, entityId := splitQual(aggEntity)
	entityShort := cmdShortForm(entityId)
	f.Commentf("Handle generically performs %s", DoSomething)
	f.Func().Params(
		Id("h").Id(DoSomething+"Handler"),
	).Id(
		"Handle",
	).Params(
		Id("ctx").Qual("context", "Context"),
		Id(cmdShortForm(DoSomething)).Id(DoSomething),
	).Params(
		Id("error"),
	).Block(
		If(
			Op("!").Id(cmdShortForm(DoSomething)).Dot("IsIdentifiable").Call(),
		).Block(
			Return().Id("Err"+DoSomething+"NotIdentifiable"),
		),
		If(
			Id("err").Op(":=").Id("h").Dot("agg").Dot(
				"Update",
			).Call(
				Id("ctx"),
				Id(cmdShortForm(DoSomething)),
				Func().Params(
					Id(entityShort).Op("*").Qual(entityImp, entityId),
				).Add(
					Id("error"),
				).Block(
					If(
						Id("ok").Op(":=").Id("h").Dot("pol").Dot(
							"Can",
						).Call(
							Id("ctx"),
							Id(cmdShortForm(DoSomething)),
							Lit(DoSomething),
							Qual("json", "Marshal").Call(Id(entityShort)),
						),
						Op("!").Id("ok"),
					).Block(
						Return().Id("ErrNotAuthorizedTo"+DoSomething),
					),
					If(
						Id("err").Op(":=").Id(cmdShortForm(DoSomething)).Dot(
							"handle",
						).Call(
							Id("ctx"),
							Id(entityShort),
						),
						Id("err").Op("!=").Id("nil"),
					).Block(
						Return().Id("err"),
					),
					Return().Id("nil"),
				),
			),
			Id("err").Op("!=").Id("nil"),
		).Block(
			Return().Id("err"),
		),
		Return().Id("nil"),
	)
}

func addCommandIsIdentifiable(f *File, DoSomething string) {
	f.Commentf("IsIdentifiable answers the question wether the %s's object is identifiable", DoSomething)
	f.Null().Func().Params(
		Id(cmdShortForm(DoSomething)).Id(DoSomething),
	).Id(
		"IsIdentifiable",
	).Params(
	).Id("bool").Block(
		If(
			Qual("reflect", "ValueOf").Call(
				Id(cmdShortForm(DoSomething)).Dot("Identifier").Call(),
			).Dot("IsZero").Call(),
		).Block(
			Return().Id("true"),
		),
		Return().Id("false"),
	)
}

// Stubs ...

func addCommandTypeStub(f *File, DoSomething string) {
	f.Commentf("%s represents a %s command", DoSomething, DoSomething)
	f.Null().Type().Id(
		DoSomething,
	).Struct(
		Id("uuid").Id("string"),
		Id("userId").Id("string"),
		Line(),
		Comment("TODO: design command event/message fields (evtl. use protobuf + protoc-gen-go)"),
	)
}

func addCommandHandleStub(f *File, DoSomething string, aggEntity string) {
	entityImp, entityId := splitQual(aggEntity)
	entityShort := cmdShortForm(entityId)
	f.Commentf("handle specifically performs %s", DoSomething)
	f.Null().Func().Params(
		Id(cmdShortForm(DoSomething)).Op("*").Id(DoSomething),
	).Id(
		"handle",
	).Params(
		Id("ctx").Qual("context", "Context"),
		Id(entityShort).Op("*").Qual(entityImp, entityId),
	).Add(
		Id("error"),
	).Block(
		Comment("TODO: implement app logic"),
		Return().Id("nil"),
	)
}

func addCommandIdentifierStub(f *File, DoSomething string) {
	f.Commentf("Identifier returns the identifier of the object on which to perform %s", DoSomething)
	f.Null().Func().Params(
		Id(cmdShortForm(DoSomething)).Id(DoSomething),
	).Id(
		"Identifier",
	).Params().Add(
		Id("string"),
	).Block(
		Return().Id(cmdShortForm(DoSomething)).Dot("uuid"),
	)
}

func addCommandUserStub(f *File, DoSomething string) {
	f.Commentf("User returns the identifier of the caller")
	f.Null().Func().Params(
		Id(cmdShortForm(DoSomething)).Id(DoSomething),
	).Id(
		"User",
	).Params().Add(
		Id("string"),
	).Block(
		Return().Id(cmdShortForm(DoSomething)).Dot("userId"),
	)
}

func addCommandElevationTokenStub(f *File, DoSomething string) {
	f.Commentf("ElevationToken returns an elevation token in posession of the caller")
	f.Null().Func().Params(
		Id(cmdShortForm(DoSomething)).Id(DoSomething),
	).Id(
		"ElevationToken",
	).Params().Add(
		Id("string"),
	).Block(
		Return().Lit(``),
	)
}

func addCommandInterfaceAssertionIdentifiable(f *File, DoSomething string, identifiable string) {
	f.Commentf("Assert that %s implements Identifiable interface!", DoSomething)
	f.Var().Id("_").Qual(splitQual(identifiable)).Op("=").Parens(Op("*").Id(DoSomething)).Call(Id("nil"))
}

func addCommandInterfaceAssertionPoliceable(f *File, DoSomething string, policeable string) {
	f.Commentf("Assert that %s implements Policeable interface!", DoSomething)
	f.Var().Id("_").Qual(splitQual(policeable)).Op("=").Parens(Op("*").Id(DoSomething)).Call(Id("nil"))
}

// Composers ...

func GenCommand(cmd, topic string, withPolicy bool, adapters []struct{ Id, Qual string }, aggEntity string) *File {
	ret := NewFile("command")
	ret.HeaderComment(fmt.Sprintf("Code generated by %s: DO NOT EDIT.", cmdGenCommand))
	ret.Line()
	ret.Commentf("Topic: %s", topic)
	ret.Line()
	addCommandIsIdentifiable(ret, cmd)
	addCommandNotAuthorizedErr(ret, cmd)
	addCommandNotIdentifiableErr(ret, cmd)
	addCommandHandlerType(ret, cmd, adapters)
	addCommandHandlerConstructor(ret, cmd, adapters)
	addCommandFuncHandle(ret, cmd, withPolicy, aggEntity)
	return ret
}

func StubCommand(cmd, topic string, withPolicy, withCommandStub bool, aggEntity, identifiable, policeable string) *File {
	ret := NewFile("command")
	ret.HeaderComment(fmt.Sprintf("Code generated by %s: THESE ARE STUBS, PLEASE EDIT.", cmdGenCommand))
	ret.Line()
	ret.Commentf("/*\n\t=== Topic: %s ===\n*/", topic)
	ret.Line()
	if withCommandStub {
		addCommandTypeStub(ret, cmd)
	}
	addCommandHandleStub(ret, cmd, aggEntity)
	ret.Line()
	ret.Comment("/*\n\t=> Identifiable interface implementation ...\n*/")
	ret.Line()
	addCommandIdentifierStub(ret, cmd)
	addCommandInterfaceAssertionIdentifiable(ret, cmd, identifiable)
	if withPolicy {
		ret.Comment("/*\n\t=> Policeable interface implementation ...\n*/")
		ret.Line()
		addCommandUserStub(ret, cmd)
		addCommandElevationTokenStub(ret, cmd)
		addCommandInterfaceAssertionPoliceable(ret, cmd, policeable)
	}
	return ret
}

func GenDoc(docFile string) *File {
	ret := NewFile("command")
	ret.PackageComment("Package command implements application layer commands")
	return ret
}

// Utils ...

func cmdShortForm(s string) string {
	re := regexp.MustCompile(`[A-Z]`)
	var b bytes.Buffer
	for _, el := range re.FindAllString(s, -1) {
		b.WriteString(strings.ToLower(el))
	}
	return b.String()
}

func splitQual(s string) (string, string) {
	imp := s[:strings.LastIndex(s, ".")]
	id := s[strings.LastIndex(s, ".")+1:]
	return imp, id
}
