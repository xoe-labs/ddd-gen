package ifaces

import (
	"context"

	"github.com/xoe-labs/ddd-gen/internal/test-svc/domain"
)

// ErrorKeeperCommandHandler handles a command in the domain and keeps domain errors
type ErrorKeeperCommandHandler interface {
	// Handle handles the command on the account entity
	Handle(ctx context.Context, a *domain.Account, ifaces ...interface{}) bool
	// Errors knows how to return domain errors
	Errors() []error
}

// FactErrorKeeperCommandHandler handles a command in the domain and keeps domain errors & facts
type FactErrorKeeperCommandHandler interface {
	// Facts knows how to return domain facts
	Facts() []interface{}
}

