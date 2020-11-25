package app

import (
	"context"
	domain "github.com/xoe-labs/ddd-gen/internal/test-svc/domain"
)

// RequiresPolicer knows to make decisions on access policy
// application requires policy adapter to implement this interface.
type RequiresPolicer interface {
	Can(ctx context.Context, p OffersPoliceable, action string, a *domain.Account) bool
}
