package app

import (
	"context"
	account "github.com/xoe-labs/ddd-gen/internal/test-svc/domain/Account"
)

// RequiresPolicer knows to make decisions on access policy
// application requires policy adapter to implement this interface.
type RequiresPolicer interface {
	Can(ctx context.Context, p OffersAuthorizable, action string, a *account.Account) bool
}
