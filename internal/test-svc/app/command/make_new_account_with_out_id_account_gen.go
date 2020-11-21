// Code generated by 'ddd-gen app command': DO NOT EDIT.

package command

import (
	"context"
	"encoding/json"
	errwrap "github.com/hashicorp/errwrap"
	gouuid "github.com/satori/go.uuid"
	error1 "github.com/xoe-labs/ddd-gen/internal/test-svc/app/error"
	policy "github.com/xoe-labs/ddd-gen/internal/test-svc/app/policy"
	repository "github.com/xoe-labs/ddd-gen/internal/test-svc/app/repository"
	account "github.com/xoe-labs/ddd-gen/internal/test-svc/domain/account"
	"reflect"
)

// Topic: Account

var (
	// ErrNotAuthorizedToMakeNewAccountWithOutId signals that the caller is not authorized to perform MakeNewAccountWithOutId
	ErrNotAuthorizedToMakeNewAccountWithOutId = error1.NewAuthorizationError("ErrNotAuthorizedToMakeNewAccountWithOutId")
	// ErrMakeNewAccountWithOutIdFailedInRepository signals that MakeNewAccountWithOutId failed in the repository layer
	ErrMakeNewAccountWithOutIdFailedInRepository = error1.NewRepositoryError("ErrMakeNewAccountWithOutIdFailedInRepository")
	// ErrMakeNewAccountWithOutIdFailedInDomain signals that MakeNewAccountWithOutId failed in the domain layer
	ErrMakeNewAccountWithOutIdFailedInDomain = error1.NewDomainError("ErrMakeNewAccountWithOutIdFailedInDomain")
)

// MakeNewAccountWithOutIdHandler knows how to perform MakeNewAccountWithOutId
type MakeNewAccountWithOutIdHandler struct {
	pol policy.Policer
	agg repository.Repository
}

// NewMakeNewAccountWithOutIdHandler returns MakeNewAccountWithOutIdHandler
func NewMakeNewAccountWithOutIdHandler(pol policy.Policer, agg repository.Repository) *MakeNewAccountWithOutIdHandler {
	if reflect.ValueOf(pol).IsZero() {
		panic("no 'pol' provided!")
	}
	if reflect.ValueOf(agg).IsZero() {
		panic("no 'agg' provided!")
	}
	return &MakeNewAccountWithOutIdHandler{pol: pol, agg: agg}
}

// Handle generically performs MakeNewAccountWithOutId
func (h MakeNewAccountWithOutIdHandler) Handle(ctx context.Context, mnawoi MakeNewAccountWithOutId) (gouuid.UUID, error) {
	var innerErr error
	var repoErr error
	identifier, repoErr := h.agg.Add(ctx, func() (a *account.Account) {
		if err := mnawoi.handle(ctx, a); err != nil {
			innerErr = errwrap.Wrap(ErrMakeNewAccountWithOutIdFailedInDomain, err)
			return nil
		}
		data, err := json.Marshal(a)
		if err != nil {
			panic(err) // invariant violation: the domain shall always be consistent!
		}
		if ok := h.pol.Can(ctx, mnawoi, "MakeNewAccountWithOutId", data); !ok {
			innerErr = ErrNotAuthorizedToMakeNewAccountWithOutId
			return nil
		}
		return a
	})
	if innerErr != nil {
		return identifier, innerErr
	}
	if repoErr != nil {
		return identifier, errwrap.Wrap(ErrMakeNewAccountWithOutIdFailedInRepository, repoErr)
	}
	return identifier, nil
}
