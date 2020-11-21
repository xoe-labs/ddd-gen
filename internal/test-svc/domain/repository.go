package domain

import (
	"context"
	"github.com/satori/go.uuid"
	account "github.com/xoe-labs/ddd-gen/internal/test-svc/domain/account"
)

type Repository interface {
	Update(ctx context.Context, i Identifiable, f func(a *account.Account) error) error
}

type Identifiable interface {
	Identifier() uuid.UUID
}
