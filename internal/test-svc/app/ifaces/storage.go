package ifaces

import (
	"context"

	"github.com/xoe-labs/ddd-gen/internal/test-svc/domain/account"
)

// Reader knows how load an account entity
type StorageReader interface {
	// Load knows how to load an account entity
	Load(ctx context.Context, target Distinguishable) (a *account.Account, err error)
}

// StorageWriteReader knows how load and persist an account entity
type StorageWriteReader interface {
	StorageReader
	// Save knows how to persist an account entity
	Save(ctx context.Context, a *account.Account) (err error)
}

// Distinguishable can be identified by both, the store and external consumers of the application
type Distinguishable interface {
	// Identifer knows how to identify an object
	// It is defined at the application layer as "interchange context"
	// between storage and external service consumer (but not the domain which acts on anonymous entties)
	Identifier() string
	// IsDistinguishablei knows how to asserts that a given Distinguishable is actually uniquely distinguishable
	IsDistinguishable() bool
}
