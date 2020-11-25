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
	StorageRWAdapterIdent = "rw"
	StorageRAdapterIdent  = "r"
	PolicyAdapterIdent    = "p"
)

func generateRequiredIfaces(genPath string, useFactStorage bool, objects *generator.Objects, adapters *generator.Adapters) error {
	// doc file
	docFile := path.Join(genPath, "doc.go")
	if fileExists(docFile) {
		if err := os.Remove(docFile); err != nil {
			return err
		}
	}
	gdf := generator.GenRequiredIfacesDoc()
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
	gsf, rTyp, rwTyp := generator.GenStorageIface(objects.Entity, useFactStorage)
	if err := gsf.Save(storageFile); err != nil {
		return err
	}
	adapters.StorageRAdapter = generator.NamedQualId{
		Name: StorageRAdapterIdent,
		QualId: generator.QualId{
			Qual: pkgPath,
			Id:   rTyp,
		},
	}
	adapters.StorageRWAdapter = generator.NamedQualId{
		Name: StorageRWAdapterIdent,
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
	gpf, typ := generator.GenIfacePolicer(objects.Entity)
	if err := gpf.Save(policyFile); err != nil {
		return err
	}
	adapters.PolicyAdapter = generator.NamedQualId{
		Name: PolicyAdapterIdent,
		QualId: generator.QualId{
			Qual: pkgPath,
			Id:   typ,
		},
	}

	// command related interfaces
	commandFile := path.Join(genPath, "command.go")
	if fileExists(commandFile) {
		if err := os.Remove(commandFile); err != nil {
			return err
		}
	}
	gcf, cmdTyp, factCmdTyp := generator.GenCmdHandlerIface(objects.Entity, useFactStorage)
	if err := gcf.Save(commandFile); err != nil {
		return err
	}
	objects.ErrorKeeperCmdHandler = generator.QualId{
		Qual: pkgPath,
		Id:   cmdTyp,
	}
	objects.FactErrorKeeperCmdHandler = generator.QualId{
		Qual: pkgPath,
		Id:   factCmdTyp,
	}

	// identity related interfaces
	identityFile := path.Join(genPath, "identity.go")
	if fileExists(identityFile) {
		if err := os.Remove(identityFile); err != nil {
			return err
		}
	}
	gif, typ := generator.GenIfaceDistinguishableAssertable()
	if err := gif.Save(identityFile); err != nil {
		return err
	}
	objects.TargetIdAssertable = generator.QualId{
		Qual: pkgPath,
		Id:   typ,
	}

	return nil
}

func generateOfferedIfaces(genPath string, objects *generator.Objects) error {
	// doc file
	docFile := path.Join(genPath, "doc.go")
	if fileExists(docFile) {
		if err := os.Remove(docFile); err != nil {
			return err
		}
	}
	gdf := generator.GenOfferedIfacesDoc()
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
	log.Printf("\t%s\n", objects.TargetIdAssertable)


	// distinguishable related interfaces
	distinguishableFile := path.Join(genPath, "distinguishable.go")
	if fileExists(distinguishableFile) {
		if err := os.Remove(distinguishableFile); err != nil {
			return err
		}
	}
	gsf, disTyp := generator.GenIfaceDistinguishable(objects)
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
	gpf, polTyp := generator.GenIfacePoliceable()
	if err := gpf.Save(policeableFile); err != nil {
		return err
	}
	objects.Actor = generator.QualId{
		Qual: pkgPath,
		Id:   polTyp,
	}

	return nil
}
