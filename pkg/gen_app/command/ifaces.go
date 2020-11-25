// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package command

import (
	"log"
	"os"
	"path"

	"golang.org/x/tools/go/packages"

	"github.com/xoe-labs/ddd-gen/pkg/gen_app/generator"
)

const (
	StorageRWIdent = "rw"
	StorageRIdent  = "r"
	PolicerIdent   = "p"
)

func generateIfaces(genPath string, useFactStorage bool, objects *generator.Objects, adapters *generator.Adapters) error {
	pkgName := "app"
	// doc file
	docFile := path.Join(genPath, "doc.go")
	if fileExists(docFile) {
		if err := os.Remove(docFile); err != nil {
			return err
		}
	}
	gdf := generator.GenAppIfacesDoc(pkgName)
	if err := gdf.Save(docFile); err != nil {
		return err
	}

	// determin the fully qualified package path
	pkgs, err := packages.Load(&packages.Config{Mode: packages.NeedName}, genPath)
	if err != nil {
		return err
	}
	pkgPath := pkgs[0].PkgPath
	log.Printf("Generating package: %s\n", pkgPath)
	log.Println("  using object interfaces ...")
	log.Printf("\t%s\n", objects.Entity)

	// storage related interfaces
	storageFile := path.Join(genPath, "storage.go")
	if fileExists(storageFile) {
		if err := os.Remove(storageFile); err != nil {
			return err
		}
	}
	gsf, rTyp, rwTyp := generator.GenStorageIface(objects.Entity, useFactStorage, pkgName)
	if err := gsf.Save(storageFile); err != nil {
		return err
	}
	adapters.StorageR = generator.NamedQualId{
		Name: StorageRIdent,
		QualId: generator.QualId{
			Qual: pkgPath,
			Id:   rTyp,
		},
	}
	adapters.StorageRW = generator.NamedQualId{
		Name: StorageRWIdent,
		QualId: generator.QualId{
			Qual: pkgPath,
			Id:   rwTyp,
		},
	}

	// policy related interfaces
	policyFile := path.Join(genPath, "policy.go")
	if fileExists(policyFile) {
		if err := os.Remove(policyFile); err != nil {
			return err
		}
	}
	gpf, typ := generator.GenIfacePolicer(objects.Entity, pkgName)
	if err := gpf.Save(policyFile); err != nil {
		return err
	}
	adapters.Policer = generator.NamedQualId{
		Name: PolicerIdent,
		QualId: generator.QualId{
			Qual: pkgPath,
			Id:   typ,
		},
	}

	// command related interfaces
	commandFile := path.Join(genPath, "domain.go")
	if fileExists(commandFile) {
		if err := os.Remove(commandFile); err != nil {
			return err
		}
	}
	gcf, cmd, fk := generator.GenCmdHandlerIface(objects.Entity, useFactStorage, pkgName)
	if err := gcf.Save(commandFile); err != nil {
		return err
	}
	if useFactStorage {
		objects.FactKeeper = generator.QualId{
			Qual: pkgPath,
			Id:   fk,
		}
	}
	objects.DomainCommandHandler = generator.QualId{
		Qual: pkgPath,
		Id:   cmd,
	}

	// identity related interfaces
	identityFile := path.Join(genPath, "identity.go")
	if fileExists(identityFile) {
		if err := os.Remove(identityFile); err != nil {
			return err
		}
	}
	gif, typ := generator.GenIfaceDistinguishableAsserter(pkgName)
	if err := gif.Save(identityFile); err != nil {
		return err
	}

	// distinguishable related interfaces
	distinguishableFile := path.Join(genPath, "distinguishable.go")
	if fileExists(distinguishableFile) {
		if err := os.Remove(distinguishableFile); err != nil {
			return err
		}
	}
	gsf, disTyp := generator.GenIfaceDistinguishable(pkgName)
	if err := gsf.Save(distinguishableFile); err != nil {
		return err
	}
	objects.Target = generator.QualId{
		Qual: pkgPath,
		Id:   disTyp,
	}

	// policeable related interfaces
	policeableFile := path.Join(genPath, "policeable.go")
	if fileExists(policeableFile) {
		if err := os.Remove(policeableFile); err != nil {
			return err
		}
	}
	gpf, polTyp := generator.GenIfacePoliceable(pkgName)
	if err := gpf.Save(policeableFile); err != nil {
		return err
	}
	objects.Actor = generator.QualId{
		Qual: pkgPath,
		Id:   polTyp,
	}
	return nil
}
