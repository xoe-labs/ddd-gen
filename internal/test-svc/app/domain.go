package app

import (
	"context"
	account "github.com/xoe-labs/ddd-gen/internal/test-svc/domain/Account"
)

// RequiresCommandHandler handles a command in the domain
type RequiresCommandHandler interface {
	// Handle handles the command on Account entity
	Handle(ctx context.Context, a *account.Account) bool
}

// RequiresErrorKeeper keeps domain errors
type RequiresErrorKeeper interface {
	// Errors knows how to return collected domain errors
	Errors() []error
}

// OffersFactKeeper keeps domain facts
type OffersFactKeeper interface {
	// Facts knows how to return domain facts
	Facts() []interface{}
}
