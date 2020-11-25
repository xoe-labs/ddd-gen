// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

// Constants represent invariant contract requirements that would
// have been too cumbersome to expose as configuration
// They _could_ be configuration, there is just not much gain in it.
const (
	// Entity
	Stringer           string = "String"
	SetterPrefix              = "Set"
	Neww                      = "New"
	MustNew                   = "MustNew"
	Equal                     = "Equal"
	UnmarshalFromStore        = "UnmarshalFromStore"
	Apply                     = "Apply"

	// DomainCommandHandler
	Handle      = "Handle"
	Facts       = "Facts"
	FactsField  = "facts"
	Errors      = "Errors"
	ErrorsField = "errors"
	RecordOn    = "recordOn"
	Raise       = "raise"
)
