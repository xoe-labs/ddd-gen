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
	StorageRAdapter    NamedQualId
	StorageRWAdapter   NamedQualId
	PolicyAdapter      NamedQualId
	DomServiceAdapters []NamedQualId
}

// Objects are represented by application level or domain interfaces
type Objects struct {
	Target                    QualId // target represents a distinguishable entity
	TargetIdAssertable        QualId // target id assertable represents a target that can be asserted to be distinguishable
	Entity                    QualId // entity represents a non-distinguishable concrete entity
	Actor                     QualId // actor represents the caller of a command
	ErrorKeeperCmdHandler     QualId // error keeper and command handler represents an object that handles the command and keeps errors
	FactErrorKeeperCmdHandler QualId // fact & error keeper and command handler represents an object that handles the command and keeps facts & errors

}

// Error constructors create error values
type Errors struct {
	AuthorizationErrorNew        QualId
	TargetIdentificationErrorNew QualId
	StorageLoadingErrorNew       QualId
	StorageSavingErrorNew        QualId
	DomainErrorNew               QualId
}
