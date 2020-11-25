package requires

import (
	"context"
	domain "github.com/xoe-labs/ddd-gen/internal/test-svc/domain"
)

// Policer knows to make decisions on access policy
// application requires policy adapter to implement this interface.
type Policer interface {
	Can(ctx context.Context, p Policeable, action string, a *domain.Account) bool
}
