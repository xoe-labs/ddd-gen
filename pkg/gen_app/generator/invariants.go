// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

// Constants represent invariant contract requirements that would
// have been too cumbersome to expose as configuration
// They _could_ be configuration, there is just not much gain in it.
const (
	TargetDistinguishableIdent            string = "OffersDistinguishable"
	TargetDistinguishableAssertableIdent         = "RequiresDistinguishableAssertable"
	TargetDistinguishableIdMethod                = "Identifier"
	TargetDistinguishableAssertMethodName        = "IsDistinguishable"

	PoliceableIdent           = "OffersPoliceable"
	PolicyAdapterIfaceIdent   = "RequiresPolicer"
	PolicyAssertionMethodName = "Can"

	StorageReaderIdent         = "RequiresStorageReader"
	StorageWriterReaderIdent   = "RequiresStorageWriterReader"
	StorageLoadMethodName      = "Load"
	StorageSaveMethodName      = "Save"
	StorageSaveFactsMethodName = "SaveFacts"

	CommandHandler                     = "commandHandler"
	CmdHandleMethodName                = "Handle"
	ErrorKeeper                        = "errorKeeper"
	ErrorKeeperCollectErrorsMethodName = "Errors"
	FactKeeper                         = "OffersFactKeeper"
	FactKeeperCollectFactsMethodName   = "Facts"
	DomainCommandHandler               = "RequiresDomainCommandHandler"
)
