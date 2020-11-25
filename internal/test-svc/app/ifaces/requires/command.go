package requires

import (
	"context"
	domain "github.com/xoe-labs/ddd-gen/internal/test-svc/domain"
)

// ErrorKeeperCommandHandler handles a command in the domain and keeps domain errors
// application requires domain to implement this interface.
type ErrorKeeperCommandHandler interface {
	// Handle handles the command on Account entity
	Handle(ctx context.Context, a *domain.Account, ifaces ...interface{}) bool
	// Errors knows how to return collected domain errors
	Errors() []error
}
