package domain

import (
	"context"
	account "github.com/xoe-labs/ddd-gen/internal/test-svc/domain/Account"
)

// Handle handles MakeNewAccount in the domain
// returns true for success or false for failure.
// record errors with mna.raise(err).
// implements application layer's CommandHandler interface.
func (mna *MakeNewAccount) Handle(ctx context.Context, a *account.Account) bool {
	return true
}
