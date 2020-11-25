// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

import (
	. "github.com/dave/jennifer/jen"
)

// Application contracts (required and offered) ...

// Required interfaces ...

func genIfaceStorageReader(f *File, entity QualId) (typIdent string) {
	entityShort := cmdShortForm(entity.Id)
	f.Commentf("%s knows how load %s entity", StorageReaderIdent, entity.Id)
	f.Comment("application requires storage adapter to implement this interface.")
	f.Type().Id(
		StorageReaderIdent,
	).Interface(
		Commentf(
			"%s knows how to load %s entity", StorageLoadMethodName, entity.Id,
		),
		Id(
			StorageLoadMethodName,
		).Params(
			Id("ctx").Qual("context", "Context"),
			Id("target").Id(
				TargetDistinguishableIdent,
			),
		).Params(
			Id(entityShort).Op("*").Qual(entity.Qual, entity.Id),
			Id("err").Id("error"),
		),
	)
	return StorageReaderIdent
}

func genIfaceStorageWriterReader(f *File, entity QualId, useFactStorage bool) (typIdent string) {
	entityShort := cmdShortForm(entity.Id)
	f.Commentf("%s knows how load and persist %s entity", StorageWriterReaderIdent, entity.Id)
	f.Comment("application requires storage adapter to implement this interface.")
	f.Type().Id(
		StorageWriterReaderIdent,
	).InterfaceFunc(func(g *Group) {
		g.Id(
			StorageReaderIdent,
		)
		if useFactStorage {
			g.Commentf(
				"%s knows how to persist domain facts on %s entity", StorageSaveFactsMethodName, entity.Id,
			)
			g.Id(
				StorageSaveFactsMethodName,
			).Params(
				Id("ctx").Qual("context", "Context"),
				Id("target").Id(
					TargetDistinguishableIdent,
				),
				Id("fk").Id(FactKeeper),
			).Params(
				Id("err").Id("error"),
			)
		} else {
			g.Commentf(
				"%s knows how to persist %s entity", StorageSaveMethodName, entity.Id,
			)
			g.Id(
				StorageSaveMethodName,
			).Params(
				Id("ctx").Qual("context", "Context"),
				Id("target").Id(
					TargetDistinguishableIdent,
				),
				Id(entityShort).Op("*").Qual(entity.Qual, entity.Id),
			).Params(
				Id("err").Id("error"),
			)
		}
	})
	return StorageWriterReaderIdent
}

func GenIfacePolicer(entity QualId) (f *File, typIdent string) {
	entityShort := cmdShortForm(entity.Id)
	f = NewFile("requires")
	f.Commentf("%s knows to make decisions on access policy", PolicyAdapterIfaceIdent)
	f.Comment("application requires policy adapter to implement this interface.")
	f.Type().Id(
		PolicyAdapterIfaceIdent,
	).Interface(
		Id("Can").Params(
			Id("ctx").Qual("context", "Context"),
			Id("p").Id(
				PoliceableIdent,
			),
			Id("action").Id("string"),
			Id(entityShort).Op("*").Qual(entity.Qual, entity.Id),
		).Params(
			Id("bool"),
		),
	)
	return f, PolicyAdapterIfaceIdent
}

func genIfaceCommandHandler(f *File, entity QualId) {
	entityShort := cmdShortForm(entity.Id)
	f.Commentf("%s handles a command in the domain", CommandHandler)
	f.Type().Id(
		CommandHandler,
	).Interface(
		Commentf(
			"%s handles the command on %s entity", CmdHandleMethodName, entity.Id,
		),
		Id(
			CmdHandleMethodName,
		).Params(
			Id("ctx").Qual("context", "Context"),
			Id(entityShort).Op("*").Qual(entity.Qual, entity.Id),
			Id("ifaces").Op("...").Interface(),
		).Params(
			Id("bool"),
		),
	)
}

func genIfaceErrorKeeper(f *File) {
	f.Commentf("%s keeps domain errors", ErrorKeeper)
	f.Type().Id(
		ErrorKeeper,
	).Interface(
		Commentf(
			"%s knows how to return collected domain errors", ErrorKeeperCollectErrorsMethodName,
		),
		Id(
			ErrorKeeperCollectErrorsMethodName,
		).Params().Params(
			Index().Id("error"),
		),
	)
}

func genIfaceFactKeeper(f *File) {
	f.Commentf("%s keeps domain facts", FactKeeper)
	f.Type().Id(
		FactKeeper,
	).Interface(
		Commentf(
			"%s knows how to return domain facts", FactKeeperCollectFactsMethodName,
		),
		Id(
			FactKeeperCollectFactsMethodName,
		).Params().Params(
			Index().Interface(),
		),
	)
}

func genIfaceDomainCommandHandler(f *File, entity QualId) (typIdent string) {
	genIfaceCommandHandler(f, entity)
	genIfaceErrorKeeper(f)
	f.Commentf("%s handles a command in the domain and keeps domain errors", DomainCommandHandler)
	f.Comment("application requires domain to implement this interface.")
	f.Type().Id(
		DomainCommandHandler,
	).Interface(
		Id(
			CommandHandler,
		),
		Id(
			ErrorKeeper,
		),
	)
	return DomainCommandHandler
}

func genIfaceDomainCommandHandlerWithFacts(f *File, entity QualId) (typIdent string) {
	genIfaceCommandHandler(f, entity)
	genIfaceErrorKeeper(f)
	genIfaceFactKeeper(f)
	f.Commentf("%s handles a command in the domain and keeps domain errors & facts", DomainCommandHandler)
	f.Comment("application requires domain to implement this interface.")
	f.Type().Id(
		DomainCommandHandler,
	).Interface(
		Id(
			CommandHandler,
		),
		Id(
			ErrorKeeper,
		),
		Id(
			FactKeeper,
		),
	)
	return DomainCommandHandler
}

func GenIfaceDistinguishableAssertable() (f *File, typIdent string) {
	f = NewFile("requires")
	f.Commentf("%sAssertable can be asserted to be distinguishable", TargetDistinguishableIdent)
	f.Commentf("application requires to be able to assert that %s can actually be identified", TargetDistinguishableIdent)
	f.Type().Id(
		TargetDistinguishableIdent+"Assertable",
	).Interface(
		Commentf("%s knows how to assert that a potential %s can be actually identified", TargetDistinguishableAssertMethodName, TargetDistinguishableIdent),
		Id(
			TargetDistinguishableAssertMethodName,
		).Params().Params(
			Id("bool"),
		),
	)
	return f, TargetDistinguishableIdent + "Assertable"
}

func GenStorageIface(entity QualId, useFactStorage bool) (f *File, storageReaderTypeIdent, storageReaderWriterTypeIdent string) {
	ret := NewFile("requires")
	storageReader := genIfaceStorageReader(ret, entity)
	storageReaderWriter := genIfaceStorageWriterReader(ret, entity, useFactStorage)
	return ret, storageReader, storageReaderWriter
}

func GenCmdHandlerIface(entity QualId, useFactStorage bool) (f *File, cmd, fk string) {
	ret := NewFile("requires")
	if useFactStorage {
		cmd = genIfaceDomainCommandHandlerWithFacts(ret, entity)
		return ret, cmd, FactKeeper
	}
	cmd = genIfaceDomainCommandHandler(ret, entity)
	return ret, cmd, ""
}

func GenRequiredIfacesDoc() *File {
	ret := NewFile("requires")
	ret.PackageComment("Package requires declares interfaces the application layer requires")
	return ret
}

// Offered interfaces ...

func GenIfaceDistinguishable(objects *Objects) (f *File, typIdent string) {
	f = NewFile("offers")
	f.Commentf("%s can be identified", TargetDistinguishableIdent)
	f.Commentf("application implements %s and thereby offers storage adapter and external consumers a common language to reason about identity", TargetDistinguishableIdent)
	f.Commentf("TODO: implement %s", TargetDistinguishableIdent)
	f.Type().Id(
		TargetDistinguishableIdent,
	).Interface(
		Qual(
			objects.TargetIdAssertable.Qual,
			objects.TargetIdAssertable.Id,
		),
		Commentf("%s knows how to identify %s", TargetDistinguishableIdMethod, TargetDistinguishableIdent),
		Comment("TODO: adapt return type to your needs "),
		Id(
			TargetDistinguishableIdMethod,
		).Params().Params(
			Id("string"),
		),
	)
	return f, TargetDistinguishableIdent
}

func GenIfacePoliceable() (f *File, typIdent string) {
	f = NewFile("offers")
	f.Commentf("%s is an actor that can be policed", PoliceableIdent)
	f.Commentf("application implements %s and thereby offers policy adapter and external consumers a common language to reason about a policeable actor", PoliceableIdent)
	f.Commentf("TODO: implement %s", PoliceableIdent)
	f.Type().Id(
		PoliceableIdent,
	).Interface(
		Comment("TODO: adapt to your needs"),
		Line(),
		Id(
			"User",
		).Params().Params(
			Id("string"),
		),
		Id(
			"ElevationToken",
		).Params().Params(
			Id("string"),
		),
	)
	return f, PoliceableIdent
}

func GenOfferedIfacesDoc() *File {
	ret := NewFile("offers")
	ret.PackageComment("Package offers declares interfaces the application offers")
	return ret
}
