// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

// Constants represent invariant contract requirements that would
// have been too cumbersome to expose as configuration
// They _could_ be configuration, there is just not much gain in it.
const (
	Distinguishable               string = "OffersDistinguishable"
	DistinguishableAsserter              = "RequiresDistinguishableAsserter"
	DistinguishableMethod                = "Identifier"
	DistinguishableAsserterMethod        = "IsDistinguishable"

	Policeable    = "OffersPoliceable"
	Policer       = "RequiresPolicer"
	PolicerMethod = "Can"

	StorageReader          = "RequiresStorageReader"
	StorageWriterReader    = "RequiresStorageWriterReader"
	StorageLoadMethod      = "Load"
	StorageSaveMethod      = "Save"
	StorageSaveFactsMethod = "SaveFacts"

	CommandHandler       = "RequiresCommandHandler"
	CommandHandlerMethod = "Handle"
	ErrorKeeper          = "RequiresErrorKeeper"
	ErrorKeeperMethod    = "Errors"
	FactKeeper           = "OffersFactKeeper"
	FactKeeperMethod     = "Facts"
)
