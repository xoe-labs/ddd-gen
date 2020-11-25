// Code generated by 'ddd-gen app command': DO NOT EDIT.

package command

import (
	"context"
	"encoding/json"
	errwrap "github.com/hashicorp/errwrap"
	error1 "github.com/xoe-labs/ddd-gen/internal/test-svc/app/error"
	policy "github.com/xoe-labs/ddd-gen/internal/test-svc/app/policy"
	repository "github.com/xoe-labs/ddd-gen/internal/test-svc/app/repository"
	account "github.com/xoe-labs/ddd-gen/internal/test-svc/domain/account"
	"reflect"
)

// Topic: Account

var (
	// ErrNotAuthorizedToMakeNewAccount signals that the caller is not authorized to perform MakeNewAccount
	ErrNotAuthorizedToMakeNewAccount = error1.NewAuthorizationError("ErrNotAuthorizedToMakeNewAccount")
	// ErrMakeNewAccountNotIdentifiable signals that MakeNewAccount's command object was not identifiable
	ErrMakeNewAccountNotIdentifiable = error1.NewIdentificationError("ErrMakeNewAccountNotIdentifiable")
	// ErrMakeNewAccountFailedInRepository signals that MakeNewAccount failed in the repository layer
	ErrMakeNewAccountFailedInRepository = error1.NewRepositoryError("ErrMakeNewAccountFailedInRepository")
	// ErrMakeNewAccountFailedInDomain signals that MakeNewAccount failed in the domain layer
	ErrMakeNewAccountFailedInDomain = error1.NewDomainError("ErrMakeNewAccountFailedInDomain")
)

// MakeNewAccountHandler knows how to perform MakeNewAccount
type MakeNewAccountHandler struct {
	pol policy.Policer
	agg repository.Repository
}

// NewMakeNewAccountHandler returns MakeNewAccountHandler
func NewMakeNewAccountHandler(pol policy.Policer, agg repository.Repository) *MakeNewAccountHandler {
	if reflect.ValueOf(pol).IsZero() {
		panic("no 'pol' provided!")
	}
	if reflect.ValueOf(agg).IsZero() {
		panic("no 'agg' provided!")
	}
	return &MakeNewAccountHandler{pol: pol, agg: agg}
}

// Handle generically performs MakeNewAccount
func (h MakeNewAccountHandler) Handle(ctx context.Context, mna MakeNewAccount) error {
	if reflect.ValueOf(mna.Identifier()).IsZero() {
		return ErrMakeNewAccountNotIdentifiable
	}
	var innerErr error
	var repoErr error
	_, repoErr := h.agg.Add(ctx, func() (a *account.Account) {
		if err := mna.handle(ctx, a); err != nil {
			innerErr = errwrap.Wrap(ErrMakeNewAccountFailedInDomain, err)
			return nil
		}
		data, err := json.Marshal(a)
		if err != nil {
			panic(err) // invariant violation: the domain shall always be consistent!
		}
		if ok := h.pol.Can(ctx, mna, "MakeNewAccount", data); !ok {
			innerErr = ErrNotAuthorizedToMakeNewAccount
			return nil
		}
		return a
	})
	if innerErr != nil {
		return innerErr
	}
	if repoErr != nil {
		return errwrap.Wrap(ErrMakeNewAccountFailedInRepository, repoErr)
	}
	return nil
}