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
	// ErrNotAuthorizedToDeleteAccount signals that the caller is not authorized to perform DeleteAccount
	ErrNotAuthorizedToDeleteAccount = error1.NewAuthorizationError("ErrNotAuthorizedToDeleteAccount")
	// ErrDeleteAccountNotIdentifiable signals that DeleteAccount's command object was not identifiable
	ErrDeleteAccountNotIdentifiable = error1.NewIdentificationError("ErrDeleteAccountNotIdentifiable")
	// ErrDeleteAccountFailedInRepository signals that DeleteAccount failed in the repository layer
	ErrDeleteAccountFailedInRepository = error1.NewRepositoryError("ErrDeleteAccountFailedInRepository")
	// ErrDeleteAccountFailedInDomain signals that DeleteAccount failed in the domain layer
	ErrDeleteAccountFailedInDomain = error1.NewDomainError("ErrDeleteAccountFailedInDomain")
)

// DeleteAccountHandler knows how to perform DeleteAccount
type DeleteAccountHandler struct {
	pol policy.Policer
	agg repository.Repository
}

// NewDeleteAccountHandler returns DeleteAccountHandler
func NewDeleteAccountHandler(pol policy.Policer, agg repository.Repository) *DeleteAccountHandler {
	if reflect.ValueOf(pol).IsZero() {
		panic("no 'pol' provided!")
	}
	if reflect.ValueOf(agg).IsZero() {
		panic("no 'agg' provided!")
	}
	return &DeleteAccountHandler{pol: pol, agg: agg}
}

// Handle generically performs DeleteAccount
func (h DeleteAccountHandler) Handle(ctx context.Context, da DeleteAccount) error {
	if reflect.ValueOf(da.Identifier()).IsZero() {
		return ErrDeleteAccountNotIdentifiable
	}
	var innerErr error
	var repoErr error
	repoErr = h.agg.Remove(ctx, da, func(a account.Account) bool {
		data, err := json.Marshal(a)
		if err != nil {
			panic(err) // invariant violation: the domain shall always be consistent!
		}
		if ok := h.pol.Can(ctx, da, "DeleteAccount", data); !ok {
			innerErr = ErrNotAuthorizedToDeleteAccount
			return false
		}
		if err := da.handle(ctx, a); err != nil {
			innerErr = errwrap.Wrap(ErrDeleteAccountFailedInDomain, err)
			return false
		}
		return true
	})
	if innerErr != nil {
		return innerErr
	}
	if repoErr != nil {
		return errwrap.Wrap(ErrDeleteAccountFailedInRepository, repoErr)
	}
	return nil
}