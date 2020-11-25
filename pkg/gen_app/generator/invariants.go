// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

// Constants represent invariant contract requirements that would
// have been too cumbersome to expose as configuration
// They _could_ be configuration, there is just not much gain in it.
const (
	TargetDistinguishableIdent            string = "Distinguishable"
	TargetDistinguishableIdMethod                = "Identifier"
	TargetDistinguishableAssertMethodName        = "IsDistinguishable"

	PoliceableIdent           = "Policeable"
	PolicyAdapterIfaceIdent   = "Policer"
	PolicyAssertionMethodName = "Can"

	StorageReaderIdent         = "StorageReader"
	StorageWriterReaderIdent   = "StorageWriterReader"
	StorageLoadMethodName      = "Load"
	StorageSaveMethodName      = "Save"
	StorageSaveFactsMethodName = "SaveFacts"

	ErrorKeeperCommandHandlerIdent        = "ErrorKeeperCommandHandler"
	FactErrorKeeperCommandHandlerIdent    = "FactErrorKeeperCommandHandler"
	CmdHandleMethodName                   = "Handle"
	ErrorKeeperCollectErrorsMethodName    = "Errors"
	FactErrorKeeperCollectFactsMethodName = "Facts"
)
