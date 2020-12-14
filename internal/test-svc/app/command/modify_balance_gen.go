// Code generated by 'ddd-gen app command': DO NOT EDIT.

package command

import (
	"context"
	errwrap "github.com/hashicorp/errwrap"
	app "github.com/xoe-labs/ddd-gen/internal/test-svc/app"
	errors "github.com/xoe-labs/ddd-gen/internal/test-svc/app/errors"
	domain "github.com/xoe-labs/ddd-gen/internal/test-svc/domain"
	"reflect"
)

// Topic: Balance

var (
	// ErrNotAuthorizedToModifyBalance signals that the caller is not authorized to perform ModifyBalance
	ErrNotAuthorizedToModifyBalance = errors.NewAuthorizationError("ErrNotAuthorizedToModifyBalance")
	// ErrModifyBalanceHasNoTarget signals that ModifyBalance's target was not distinguishable
	ErrModifyBalanceHasNoTarget = errors.NewTargetIdentificationError("ErrModifyBalanceHasNoTarget")
	// ErrModifyBalanceLoadingFailed signals that ModifyBalance storage failed to load the entity
	ErrModifyBalanceLoadingFailed = errors.NewStorageLoadingError("ErrModifyBalanceLoadingFailed")
	// ErrModifyBalanceSavingFailed signals that ModifyBalance failed to save the entity
	ErrModifyBalanceSavingFailed = errors.NewStorageSavingError("ErrModifyBalanceSavingFailed")
	// ErrModifyBalanceFailedInDomain signals that ModifyBalance failed in the domain layer
	ErrModifyBalanceFailedInDomain = errors.NewDomainError("ErrModifyBalanceFailedInDomain")
)

// ModifyBalanceHandlerWrapper knows how to perform ModifyBalance
type ModifyBalanceHandlerWrapper struct {
	rw app.RequiresStorageWriterReader
	p  app.RequiresPolicer
}

// NewModifyBalanceHandlerWrapper returns ModifyBalanceHandlerWrapper
func NewModifyBalanceHandlerWrapper(rw app.RequiresStorageWriterReader, p app.RequiresPolicer) *ModifyBalanceHandlerWrapper {
	if reflect.ValueOf(rw).IsZero() {
		panic("no 'rw' provided!")
	}
	if reflect.ValueOf(p).IsZero() {
		panic("no 'p' provided!")
	}
	return &ModifyBalanceHandlerWrapper{rw: rw, p: p}
}

// Handle generically performs ModifyBalance
func (h ModifyBalanceHandlerWrapper) Handle(ctx context.Context, mb domain.ModifyBalance, actor app.OffersAuthorizable, target app.OffersDistinguishable) error {
	// assert that target is distinguishable
	if !target.IsDistinguishable() {
		return ErrModifyBalanceHasNoTarget
	}
	// load entity from store; handle + wrap error
	a, loadErr := h.rw.Load(ctx, target)
	if loadErr != nil {
		return errwrap.Wrap(ErrModifyBalanceLoadingFailed, loadErr)
	}
	// assert authorization via policy interface
	if ok := h.p.Can(ctx, actor, "ModifyBalance", a); !ok {
		// return opaque error: handle potentially sensitive policy errors out-of-band!
		return ErrNotAuthorizedToModifyBalance
	}
	// assert correct command handling by the domain
	if ok := mb.Handle(ctx, a); !ok {
		var domErr error
		for i, e := range mb.Errors() {
			if i == 0 {
				domErr = e
			} else {
				domErr = errwrap.Wrap(domErr, e)
			}
		}
		return ErrModifyBalanceFailedInDomain
	}
	// save domain facts to storage
	saveErr := h.rw.SaveFacts(ctx, target, app.OffersFactKeeper(&mb))
	if saveErr != nil {
		return errwrap.Wrap(ErrModifyBalanceSavingFailed, saveErr)
	}
	return nil
}

// compile time assertions
var (
	_ app.RequiresCommandHandler = (*domain.ModifyBalance)(nil)
	_ app.RequiresErrorKeeper    = (*domain.ModifyBalance)(nil)
	_ app.OffersFactKeeper       = (*domain.ModifyBalance)(nil)
)