package requires

import (
	"context"
	domain "github.com/xoe-labs/ddd-gen/internal/test-svc/domain"
)

// StorageReader knows how load Account entity
// application requires storage adapter to implement this interface.
type StorageReader interface {
	// Load knows how to load Account entity
	Load(ctx context.Context, target Distinguishable) (a *domain.Account, err error)
}

// StorageWriterReader knows how load and persist Account entity
// application requires storage adapter to implement this interface.
type StorageWriterReader interface {
	StorageReader
	// SaveFacts knows how to persist domain facts on Account entity
	SaveFacts(ctx context.Context, target Distinguishable, fk FactKeeper) (err error)
}
