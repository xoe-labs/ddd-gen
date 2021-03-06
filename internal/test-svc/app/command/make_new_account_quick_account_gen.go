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

// Topic: Account

var (
	// ErrNotAuthorizedToMakeNewAccountQuick signals that the caller is not authorized to perform MakeNewAccountQuick
	ErrNotAuthorizedToMakeNewAccountQuick = errors.NewAuthorizationError("ErrNotAuthorizedToMakeNewAccountQuick")
	// ErrMakeNewAccountQuickHasNoTarget signals that MakeNewAccountQuick's target was not distinguishable
	ErrMakeNewAccountQuickHasNoTarget = errors.NewTargetIdentificationError("ErrMakeNewAccountQuickHasNoTarget")
	// ErrMakeNewAccountQuickLoadingFailed signals that MakeNewAccountQuick storage failed to load the entity
	ErrMakeNewAccountQuickLoadingFailed = errors.NewStorageLoadingError("ErrMakeNewAccountQuickLoadingFailed")
	// ErrMakeNewAccountQuickSavingFailed signals that MakeNewAccountQuick failed to save the entity
	ErrMakeNewAccountQuickSavingFailed = errors.NewStorageSavingError("ErrMakeNewAccountQuickSavingFailed")
	// ErrMakeNewAccountQuickFailedInDomain signals that MakeNewAccountQuick failed in the domain layer
	ErrMakeNewAccountQuickFailedInDomain = errors.NewDomainError("ErrMakeNewAccountQuickFailedInDomain")
)

// MakeNewAccountQuickHandlerWrapper knows how to perform MakeNewAccountQuick
type MakeNewAccountQuickHandlerWrapper struct {
	rw app.RequiresStorageWriterReader
	p  app.RequiresPolicer
}

// NewMakeNewAccountQuickHandlerWrapper returns MakeNewAccountQuickHandlerWrapper
func NewMakeNewAccountQuickHandlerWrapper(rw app.RequiresStorageWriterReader, p app.RequiresPolicer) *MakeNewAccountQuickHandlerWrapper {
	if reflect.ValueOf(rw).IsZero() {
		panic("no 'rw' provided!")
	}
	if reflect.ValueOf(p).IsZero() {
		panic("no 'p' provided!")
	}
	return &MakeNewAccountQuickHandlerWrapper{rw: rw, p: p}
}

// Handle generically performs MakeNewAccountQuick
func (h MakeNewAccountQuickHandlerWrapper) Handle(ctx context.Context, mnaq domain.MakeNewAccountQuick, actor app.OffersAuthorizable, target app.OffersDistinguishable) error {
	// assert that target is distinguishable
	if !target.IsDistinguishable() {
		return ErrMakeNewAccountQuickHasNoTarget
	}
	// load entity from store; handle + wrap error
	a, loadErr := h.rw.Load(ctx, target)
	if loadErr != nil {
		return errwrap.Wrap(ErrMakeNewAccountQuickLoadingFailed, loadErr)
	}
	// assert authorization via policy interface
	if ok := h.p.Can(ctx, actor, "MakeNewAccountQuick", a); !ok {
		// return opaque error: handle potentially sensitive policy errors out-of-band!
		return ErrNotAuthorizedToMakeNewAccountQuick
	}
	// assert correct command handling by the domain
	if ok := mnaq.Handle(ctx, a); !ok {
		var domErr error
		for i, e := range mnaq.Errors() {
			if i == 0 {
				domErr = e
			} else {
				domErr = errwrap.Wrap(domErr, e)
			}
		}
		return ErrMakeNewAccountQuickFailedInDomain
	}
	// save domain facts to storage
	saveErr := h.rw.SaveFacts(ctx, target, app.OffersFactKeeper(&mnaq))
	if saveErr != nil {
		return errwrap.Wrap(ErrMakeNewAccountQuickSavingFailed, saveErr)
	}
	return nil
}

// compile time assertions
var (
	_ app.RequiresCommandHandler = (*domain.MakeNewAccountQuick)(nil)
	_ app.RequiresErrorKeeper    = (*domain.MakeNewAccountQuick)(nil)
	_ app.OffersFactKeeper       = (*domain.MakeNewAccountQuick)(nil)
)
