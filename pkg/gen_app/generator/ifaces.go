// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

import (
	"fmt"
	. "github.com/dave/jennifer/jen"
)

// Application contracts (required and offered) ...

func GenAppIfacesDoc(pkgName string) *File {
	ret := NewFile(pkgName)
	ret.PackageComment(
		fmt.Sprintf(
			"Package %s declares interfaces which the application layer either requires or offers.",
			pkgName,
		),
	)
	ret.PackageComment("")
	ret.PackageComment("By convention, the following prefixes further qualify the interfaces:")
	ret.PackageComment("\tOffers*")
	ret.PackageComment("\tRequires*")
	ret.PackageComment("")
	ret.PackageComment("The name of the go file (e.g. `storage.go`) signifies the adapter or object of the interface.")
	ret.PackageComment("")
	ret.PackageComment(
		fmt.Sprintf("Names terminating in 'able' represent types for which %s offers an implementation:", pkgName),
	)
	ret.PackageComment("Adapters or ports shall understand those interface types as common language comming from external services")
	ret.PackageComment("Hence, their implementation is part of the package's public api.")
	return ret
}

// Required interfaces ...

func genIfaceStorageReader(f *File, entity QualId) (typIdent string) {
	entityShort := cmdShortForm(entity.Id)
	f.Commentf("%s knows how load %s entity", StorageReader, entity.Id)
	f.Comment("application requires storage adapter to implement this interface.")
	f.Type().Id(
		StorageReader,
	).Interface(
		Commentf(
			"%s knows how to load %s entity", StorageLoadMethod, entity.Id,
		),
		Id(
			StorageLoadMethod,
		).Params(
			Id("ctx").Qual("context", "Context"),
			Id("target").Id(
				Distinguishable,
			),
		).Params(
			Id(entityShort).Op("*").Qual(entity.Qual, entity.Id),
			Id("err").Id("error"),
		),
	)
	return StorageReader
}

func genIfaceStorageWriterReader(f *File, entity QualId, useFactStorage bool) (typIdent string) {
	entityShort := cmdShortForm(entity.Id)
	f.Commentf("%s knows how load and persist %s entity", StorageWriterReader, entity.Id)
	f.Comment("application requires storage adapter to implement this interface.")
	f.Type().Id(
		StorageWriterReader,
	).InterfaceFunc(func(g *Group) {
		g.Id(
			StorageReader,
		)
		if useFactStorage {
			g.Commentf(
				"%s knows how to persist domain facts on %s entity", StorageSaveFactsMethod, entity.Id,
			)
			g.Id(
				StorageSaveFactsMethod,
			).Params(
				Id("ctx").Qual("context", "Context"),
				Id("target").Id(
					Distinguishable,
				),
				Id("fk").Id(FactKeeper),
			).Params(
				Id("err").Id("error"),
			)
		} else {
			g.Commentf(
				"%s knows how to persist %s entity", StorageSaveMethod, entity.Id,
			)
			g.Id(
				StorageSaveMethod,
			).Params(
				Id("ctx").Qual("context", "Context"),
				Id("target").Id(
					Distinguishable,
				),
				Id(entityShort).Op("*").Qual(entity.Qual, entity.Id),
			).Params(
				Id("err").Id("error"),
			)
		}
	})
	return StorageWriterReader
}

func GenIfacePolicer(entity QualId, pkgName string) (f *File, typIdent string) {
	entityShort := cmdShortForm(entity.Id)
	f = NewFile(pkgName)
	f.Commentf("%s knows to make decisions on access policy", Policer)
	f.Comment("application requires policy adapter to implement this interface.")
	f.Type().Id(
		Policer,
	).Interface(
		Id("Can").Params(
			Id("ctx").Qual("context", "Context"),
			Id("p").Id(
				Policeable,
			),
			Id("action").Id("string"),
			Id(entityShort).Op("*").Qual(entity.Qual, entity.Id),
		).Params(
			Id("bool"),
		),
	)
	return f, Policer
}

func genIfaceCommandHandler(f *File, entity QualId) (typIdent string) {
	entityShort := cmdShortForm(entity.Id)
	f.Commentf("%s handles a command in the domain", CommandHandler)
	f.Type().Id(
		CommandHandler,
	).Interface(
		Commentf(
			"%s handles the command on %s entity", CommandHandlerMethod, entity.Id,
		),
		Id(
			CommandHandlerMethod,
		).Params(
			Id("ctx").Qual("context", "Context"),
			Id(entityShort).Op("*").Qual(entity.Qual, entity.Id),
			// Id("ifaces").Op("...").Interface(),
		).Params(
			Id("bool"),
		),
	)
	return CommandHandler
}

func genIfaceErrorKeeper(f *File) (typIdent string) {
	f.Commentf("%s keeps domain errors", ErrorKeeper)
	f.Type().Id(
		ErrorKeeper,
	).Interface(
		Commentf(
			"%s knows how to return collected domain errors", ErrorKeeperMethod,
		),
		Id(
			ErrorKeeperMethod,
		).Params().Params(
			Index().Id("error"),
		),
	)
	return ErrorKeeper
}

func genIfaceFactKeeper(f *File) (typIdent string) {
	f.Commentf("%s keeps domain facts", FactKeeper)
	f.Type().Id(
		FactKeeper,
	).Interface(
		Commentf(
			"%s knows how to return domain facts", FactKeeperMethod,
		),
		Id(
			FactKeeperMethod,
		).Params().Params(
			Index().Interface(),
		),
	)
	return FactKeeper
}

func GenIfaceDistinguishableAsserter(pkgName string) (f *File, typIdent string) {
	f = NewFile(pkgName)
	f.Commentf("%s can be asserted to be distinguishable", DistinguishableAsserter)
	f.Commentf("application requires to be able to assert that %s can actually be identified", Distinguishable)
	f.Type().Id(
		DistinguishableAsserter,
	).Interface(
		Commentf("%s knows how to assert that a potential %s can be actually identified", DistinguishableAsserterMethod, Distinguishable),
		Id(
			DistinguishableAsserterMethod,
		).Params().Params(
			Id("bool"),
		),
	)
	return f, DistinguishableAsserter
}

func GenStorageIface(entity QualId, useFactStorage bool, pkgName string) (f *File, storageReaderTypeIdent, storageReaderWriterTypeIdent string) {
	ret := NewFile(pkgName)
	storageReader := genIfaceStorageReader(ret, entity)
	storageReaderWriter := genIfaceStorageWriterReader(ret, entity, useFactStorage)
	return ret, storageReader, storageReaderWriter
}

func GenCmdHandlerIface(entity QualId, useFactStorage bool, pkgName string) (f *File, cmd, ek, fk string) {
	ret := NewFile(pkgName)
	cmd = genIfaceCommandHandler(ret, entity)
	ek = genIfaceErrorKeeper(ret)
	if useFactStorage {
		fk = genIfaceFactKeeper(ret)
		return ret, cmd, ek, fk
	}
	return ret, cmd, ek, ""
}

// Offered interfaces ...

func GenIfaceDistinguishable(pkgName string) (f *File, typIdent string) {
	f = NewFile(pkgName)
	f.Commentf("%s can be identified", Distinguishable)
	f.Commentf("application implements %s and thereby offers storage adapter and external consumers a common language to reason about identity", Distinguishable)
	f.Commentf("TODO: implement %s", Distinguishable)
	f.Type().Id(
		Distinguishable,
	).Interface(
		Id(DistinguishableAsserter),
		Commentf("%s knows how to identify %s", DistinguishableMethod, Distinguishable),
		Comment("TODO: adapt return type to your needs "),
		Id(
			DistinguishableMethod,
		).Params().Params(
			Id("string"),
		),
	)
	return f, Distinguishable
}

func GenIfacePoliceable(pkgName string) (f *File, typIdent string) {
	f = NewFile(pkgName)
	f.Commentf("%s is an actor that can be policed", Policeable)
	f.Commentf("application implements %s and thereby offers policy adapter and external consumers a common language to reason about a policeable actor", Policeable)
	f.Commentf("TODO: implement %s", Policeable)
	f.Type().Id(
		Policeable,
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
	return f, Policeable
}
