package app

import (
	"context"
	domain "github.com/xoe-labs/ddd-gen/internal/test-svc/domain"
)

// commandHandler handles a command in the domain
type commandHandler interface {
	// Handle handles the command on Account entity
	Handle(ctx context.Context, a *domain.Account, ifaces ...interface{}) bool
}

// errorKeeper keeps domain errors
type errorKeeper interface {
	// Errors knows how to return collected domain errors
	Errors() []error
}

// OffersFactKeeper keeps domain facts
type OffersFactKeeper interface {
	// Facts knows how to return domain facts
	Facts() []interface{}
}

// RequiresDomainCommandHandler handles a command in the domain and keeps domain errors & facts
// application requires domain to implement this interface.
type RequiresDomainCommandHandler interface {
	commandHandler
	errorKeeper
	OffersFactKeeper
}
