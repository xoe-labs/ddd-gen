// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
	// "log"
)

var cmdGenCommand string = "ddd-gen app command"

// CommandHandlerWrapper ...

func addCommandHandlerWrapperErrors(f *File,
	DoSomething string,
	assertAuthorization bool,
	errors Errors) {
	f.Null().Var().DefsFunc(func(g *Group) {
		if assertAuthorization {
			g.Commentf("ErrNotAuthorizedTo%s signals that the caller is not authorized to perform %s", DoSomething, DoSomething)
			g.Id("ErrNotAuthorizedTo"+DoSomething).Op("=").Qual(
				errors.AuthorizationErrorNew.Qual,
				errors.AuthorizationErrorNew.Id,
			).Call(
				Lit("ErrNotAuthorizedTo" + DoSomething),
			)
		}
		g.Commentf("Err%sHasNoTarget signals that %s's target was not distinguishable", DoSomething, DoSomething)
		g.Id("Err"+DoSomething+"HasNoTarget").Op("=").Qual(
			errors.TargetIdentificationErrorNew.Qual,
			errors.TargetIdentificationErrorNew.Id,
		).Call(
			Lit("Err" + DoSomething + "HasNoTarget"),
		)
		g.Commentf("Err%sLoadingFailed signals that %s storage failed to load the entity", DoSomething, DoSomething)
		g.Id("Err"+DoSomething+"LoadingFailed").Op("=").Qual(
			errors.StorageLoadingErrorNew.Qual,
			errors.StorageLoadingErrorNew.Id,
		).Call(
			Lit("Err" + DoSomething + "LoadingFailed"),
		)
		g.Commentf("Err%sSavingFailed signals that %s failed to save the entity", DoSomething, DoSomething)
		g.Id("Err"+DoSomething+"SavingFailed").Op("=").Qual(
			errors.StorageSavingErrorNew.Qual,
			errors.StorageSavingErrorNew.Id,
		).Call(
			Lit("Err" + DoSomething + "SavingFailed"),
		)
		g.Commentf("Err%sFailedInDomain signals that %s failed in the domain layer", DoSomething, DoSomething)
		g.Id("Err"+DoSomething+"FailedInDomain").Op("=").Qual(
			errors.DomainErrorNew.Qual,
			errors.DomainErrorNew.Id,
		).Call(
			Lit("Err" + DoSomething + "FailedInDomain"),
		)
	})
}

func addCommandHandlerWrapperType(f *File,
	DoSomething string,
	assertAuthorization bool,
	adapters Adapters) {
	f.Commentf("%sHandlerWrapper knows how to perform %s", DoSomething, DoSomething)
	f.Null().Type().Id(
		DoSomething + "HandlerWrapper",
	).StructFunc(func(g *Group) {
		g.Id(adapters.StorageRWAdapter.Name).Qual(adapters.StorageRWAdapter.Qual, adapters.StorageRWAdapter.Id)
		if assertAuthorization {
			g.Id(adapters.PolicyAdapter.Name).Qual(adapters.PolicyAdapter.Qual, adapters.PolicyAdapter.Id)
		}
		for _, a := range adapters.DomServiceAdapters {
			g.Id(a.Name).Qual(a.Qual, a.Id)
		}
	})
}
func addCommandHandlerWrapperConstructor(f *File,
	DoSomething string,
	assertAuthorization bool,
	adapters Adapters) {
	usedAdapters := append(adapters.DomServiceAdapters, adapters.StorageRWAdapter)
	if assertAuthorization {
		usedAdapters = append(usedAdapters, adapters.PolicyAdapter)
	}
	f.Commentf("New%sHandlerWrapper returns %sHandlerWrapper", DoSomething, DoSomething)
	f.Func().Id(
		"New" + DoSomething + "HandlerWrapper",
	).ParamsFunc(func(g *Group) {
		for _, a := range usedAdapters {
			g.Id(a.Name).Qual(a.Qual, a.Id)
		}
	}).Params(
		Op("*").Id(DoSomething + "HandlerWrapper"),
	).BlockFunc(func(g *Group) {
		for _, a := range usedAdapters {
			g.If(
				Qual(
					"reflect", "ValueOf",
				).Call(
					Id(a.Name),
				).Dot(
					"IsZero",
				).Call(),
			).Block(
				Id("panic").Call(Lit("no '" + a.Name + "' provided!")),
			)
		}
		g.Return().Op("&").Id(
			DoSomething + "HandlerWrapper",
		).ValuesFunc(func(g *Group) {
			for _, a := range usedAdapters {
				g.Id(a.Name).Op(":").Id(a.Name)
			}
		})
	})
}

func addCommandFuncHandle(f *File,
	DoSomething string,
	assertAuthorization,
	useFactStorage bool,
	objects Objects,
	adapters Adapters) {
	entityShort := cmdShortForm(objects.Entity.Id)
	f.Commentf("Handle generically performs %s", DoSomething)
	f.Func().Params(
		Id("h").Id(DoSomething + "HandlerWrapper"),
	).Id(
		"Handle",
	).ParamsFunc(func(g *Group) {
		g.Id("ctx").Qual("context", "Context")
		if useFactStorage {
			g.Id(cmdShortForm(DoSomething)).Qual(objects.FactErrorKeeperCmdHandler.Qual, objects.ErrorKeeperCmdHandler.Id)
		} else {
			g.Id(cmdShortForm(DoSomething)).Qual(objects.ErrorKeeperCmdHandler.Qual, objects.ErrorKeeperCmdHandler.Id)
		}
		g.Id("actor").Qual(objects.Actor.Qual, objects.Actor.Id)
		g.Id("target").Qual(objects.Target.Qual, objects.Target.Id)
	}).Parens(
		List(
			Id("error"),
		),
	).BlockFunc(func(g *Group) {
		g.Comment("assert that target is distinguishable")
		g.If(
			Op("!").Id("target").Dot(TargetDistinguishableAssertMethodName).Call(),
		).Block(
			Return().Id(
				"Err" + DoSomething + "HasNoTarget",
			),
		)

		g.Comment("load entity from store; handle + wrap error")
		g.List(
			Id(entityShort),
			Id("loadErr"),
		).Op(":=").Id("h").Dot(adapters.StorageRWAdapter.Name).Dot(
			StorageLoadMethodName,
		).Call(
			Id("ctx"),
			Id("target"),
		)
		g.If(
			Id("loadErr").Op("!=").Id("nil"),
		).Block(
			Return().Qual(
				"github.com/hashicorp/errwrap",
				"Wrap",
			).Call(
				Id("Err"+DoSomething+"LoadingFailed"),
				Id("loadErr"),
			),
		)

		if assertAuthorization {
			g.Comment("assert authorization via policy interface")
			g.If(
				Id(
					"ok",
				).Op(":=").Id("h").Dot(adapters.PolicyAdapter.Name).Dot(
					PolicyAssertionMethodName,
				).Call(
					Id("ctx"),
					Id("actor"),
					Lit(DoSomething),
					Id(entityShort),
				),
				Op("!").Id("ok"),
			).Block(
				Comment("return opaque error: handle potentially sensitive policy errors out-of-band!"),
				Return().Id(
					"ErrNotAuthorizedTo"+DoSomething,
				),
			)
		}

		g.Comment("assert correct command handling by the domain")
		g.If(
			Id("ok").Op(":=").Id(
				cmdShortForm(DoSomething),
			).Dot(
				CmdHandleMethodName,
			).CallFunc(func(g *Group) {
				g.Id("ctx")
				g.Id(entityShort)
				for _, a := range adapters.DomServiceAdapters {
					g.Op("&").Id("h").Dot(a.Name)
				}
			}),
			Op("!").Id("ok"),
		).Block(
			Var().Id("domErr").Id("error"),
			Commentf("%s is an ErrorKeeper", cmdShortForm(DoSomething)),
			For(
				List(
					Id("i"),
					Id("e"),
				).Op(":=").Range().Id(
					cmdShortForm(DoSomething),
				).Dot(
					ErrorKeeperCollectErrorsMethodName,
				).Call(),
			).Block(
				If(Id("i")).Op("==").Lit(0).Block(
					Id("domErr").Op("=").Id("domErr"),
				).Else().Block(
					Id("domErr").Op("=").Qual(
						"github.com/hashicorp/errwrap",
						"Wrap",
					).Call(
						Id("domErr"),
						Id("e"),
					),
				),
			),
			Return().Id(
				"Err"+DoSomething+"FailedInDomain",
			),
		)

		if useFactStorage { // a event sourcing storage
			g.Comment("save domain facts to storage")
			g.Commentf("%s is a ErrorFactKeeper", cmdShortForm(DoSomething))
			g.Id(
				"saveErr",
			).Op(":=").Id("h").Dot(adapters.StorageRWAdapter.Name).Dot(
				StorageSaveFactsMethodName,
			).Call(
				Id("ctx"),
				Id("target"),
				Id(cmdShortForm(DoSomething)),
			)
			g.If(
				Id("saveErr").Op("!=").Id("nil"),
			).Block(
				Return().Qual(
					"github.com/hashicorp/errwrap",
					"Wrap",
				).Call(
					Id("Err"+DoSomething+"SavingFailed"),
					Id("saveErr"),
				),
			)
		} else { // a modelStorage
			g.Comment("save entity to storage")
			g.Id(
				"saveErr",
			).Op(":=").Id("h").Dot(adapters.StorageRWAdapter.Name).Dot(
				StorageSaveMethodName,
			).Call(
				Id("ctx"),
				Id("target"),
				Id(entityShort),
			)
			g.If(
				Id("saveErr").Op("!=").Id("nil"),
			).Block(
				Return().Qual(
					"github.com/hashicorp/errwrap",
					"Wrap",
				).Call(
					Id("Err"+DoSomething+"SavingFailed"),
					Id("saveErr"),
				),
			)
		}
		g.Return().Id("nil")

	})
}

// Composers ...

func GenCommandHandlerWrapper(cmd,
	topic string,
	useFactStorage,
	withPolicyEnforcement bool,
	adapters Adapters,
	objects Objects,
	errors Errors) *File {
	ret := NewFile("command")
	ret.HeaderComment(fmt.Sprintf("Code generated by '%s': DO NOT EDIT.", cmdGenCommand))
	ret.Line()
	ret.Commentf("Topic: %s", topic)
	ret.Line()
	addCommandHandlerWrapperErrors(ret, cmd,
		withPolicyEnforcement,
		errors)
	addCommandHandlerWrapperType(ret, cmd,
		withPolicyEnforcement,
		adapters)
	addCommandHandlerWrapperConstructor(ret, cmd,
		withPolicyEnforcement,
		adapters)
	addCommandFuncHandle(ret, cmd,
		withPolicyEnforcement,
		useFactStorage,
		objects,
		adapters)
	return ret
}

func GenCommandDoc(docFile string) *File {
	ret := NewFile("command")
	ret.PackageComment("Package command implements application layer command wrappers")
	return ret
}