package repository

import (
	"context"
	"github.com/satori/go.uuid"
	account "github.com/xoe-labs/ddd-gen/internal/test-svc/domain/account"
)

// Repository knows how to persist the service's aggregate
type Repository interface {
	// Add knows how to create an initialized instance of an aggregate
	// it expects an initialized instance be return from an add funcion or nil to bail out
	Add(ctx context.Context,f func() (*account.Account)) (uuid.UUID, error)
	// Add knows how to remove an identifiable instance of an aggregate
	// it also returns a copy to a remove funtion in order to bail out
	Rem(ctx context.Context, i Identifiable, f func(a account.Account) bool) error
	// Update knows how to update an identifiable instance of an aggregate
	Update(ctx context.Context, i Identifiable, f func(a *account.Account) bool) error
}

// Identifiable can be identified by the Repository
type Identifiable interface {
	// Identifer knows how to identify an object
	Identifier() uuid.UUID

	// IsIdentifiable answers the question wether an object is identified
	IsIdentifiable() bool
}
