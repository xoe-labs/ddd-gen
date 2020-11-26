// Copyright © 2020 David Arnold <dar@xoe.solutions>
// SPDX-License-Identifier: MIT

package generator

type QualId struct{ Id, Qual string }

type NamedQualId struct {
	Name string
	QualId
}

// Adapters provide interfaces to the outer world
type Adapters struct {
	StorageR           NamedQualId
	StorageRW          NamedQualId
	Policer            NamedQualId
	DomServiceAdapters []NamedQualId
}

// Objects are represented by application level or domain interfaces
type Objects struct {
	Target         QualId // target represents a distinguishable entity
	Entity         QualId // entity represents a non-distinguishable concrete entity
	Actor          QualId // actor represents the caller of a command
	CommandHandler QualId // command handler handles domain commands
	ErrorKeeper    QualId // error keeper keeps domain errors
	FactKeeper     QualId // fact keeper keeps domain facts
	Domain         QualId // a qual only referncinf the domain import path

}

// Error constructors create error values
type Errors struct {
	AuthorizationErrorNew        QualId
	TargetIdentificationErrorNew QualId
	StorageLoadingErrorNew       QualId
	StorageSavingErrorNew        QualId
	DomainErrorNew               QualId
}
