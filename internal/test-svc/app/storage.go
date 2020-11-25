package app

import (
	"context"
	account "github.com/xoe-labs/ddd-gen/internal/test-svc/domain/Account"
)

// RequiresStorageReader knows how load Account entity
// application requires storage adapter to implement this interface.
type RequiresStorageReader interface {
	// Load knows how to load Account entity
	Load(ctx context.Context, target OffersDistinguishable) (a *account.Account, err error)
}

// RequiresStorageWriterReader knows how load and persist Account entity
// application requires storage adapter to implement this interface.
type RequiresStorageWriterReader interface {
	RequiresStorageReader
	// SaveFacts knows how to persist domain facts on Account entity
	SaveFacts(ctx context.Context, target OffersDistinguishable, fk OffersFactKeeper) (err error)
}
