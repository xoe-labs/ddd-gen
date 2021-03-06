// Code generated by 'ddd-gen app command': DO NOT EDIT.

package command

import (
	"context"
	errwrap "github.com/hashicorp/errwrap"
	app "github.com/xoe-labs/ddd-gen/internal/test-svc/app"
	errors "github.com/xoe-labs/ddd-gen/internal/test-svc/app/errors"
	ifaces "github.com/xoe-labs/ddd-gen/internal/test-svc/app/ifaces"
	domain "github.com/xoe-labs/ddd-gen/internal/test-svc/domain"
	"reflect"
)

// Topic: Balance

var (
	// ErrNotAuthorizedToModifyBalanceFromSvc signals that the caller is not authorized to perform ModifyBalanceFromSvc
	ErrNotAuthorizedToModifyBalanceFromSvc = errors.NewAuthorizationError("ErrNotAuthorizedToModifyBalanceFromSvc")
	// ErrModifyBalanceFromSvcHasNoTarget signals that ModifyBalanceFromSvc's target was not distinguishable
	ErrModifyBalanceFromSvcHasNoTarget = errors.NewTargetIdentificationError("ErrModifyBalanceFromSvcHasNoTarget")
	// ErrModifyBalanceFromSvcLoadingFailed signals that ModifyBalanceFromSvc storage failed to load the entity
	ErrModifyBalanceFromSvcLoadingFailed = errors.NewStorageLoadingError("ErrModifyBalanceFromSvcLoadingFailed")
	// ErrModifyBalanceFromSvcSavingFailed signals that ModifyBalanceFromSvc failed to save the entity
	ErrModifyBalanceFromSvcSavingFailed = errors.NewStorageSavingError("ErrModifyBalanceFromSvcSavingFailed")
	// ErrModifyBalanceFromSvcFailedInDomain signals that ModifyBalanceFromSvc failed in the domain layer
	ErrModifyBalanceFromSvcFailedInDomain = errors.NewDomainError("ErrModifyBalanceFromSvcFailedInDomain")
)

// ModifyBalanceFromSvcHandlerWrapper knows how to perform ModifyBalanceFromSvc
type ModifyBalanceFromSvcHandlerWrapper struct {
	rw  app.RequiresStorageWriterReader
	p   app.RequiresPolicer
	svc ifaces.Balancer
}

// NewModifyBalanceFromSvcHandlerWrapper returns ModifyBalanceFromSvcHandlerWrapper
func NewModifyBalanceFromSvcHandlerWrapper(svc ifaces.Balancer, rw app.RequiresStorageWriterReader, p app.RequiresPolicer) *ModifyBalanceFromSvcHandlerWrapper {
	if reflect.ValueOf(svc).IsZero() {
		panic("no 'svc' provided!")
	}
	if reflect.ValueOf(rw).IsZero() {
		panic("no 'rw' provided!")
	}
	if reflect.ValueOf(p).IsZero() {
		panic("no 'p' provided!")
	}
	return &ModifyBalanceFromSvcHandlerWrapper{svc: svc, rw: rw, p: p}
}

// Handle generically performs ModifyBalanceFromSvc
func (h ModifyBalanceFromSvcHandlerWrapper) Handle(ctx context.Context, mbfs domain.ModifyBalanceFromSvc, actor app.OffersAuthorizable, target app.OffersDistinguishable) error {
	// assert that target is distinguishable
	if !target.IsDistinguishable() {
		return ErrModifyBalanceFromSvcHasNoTarget
	}
	// load entity from store; handle + wrap error
	a, loadErr := h.rw.Load(ctx, target)
	if loadErr != nil {
		return errwrap.Wrap(ErrModifyBalanceFromSvcLoadingFailed, loadErr)
	}
	// assert authorization via policy interface
	if ok := h.p.Can(ctx, actor, "ModifyBalanceFromSvc", a); !ok {
		// return opaque error: handle potentially sensitive policy errors out-of-band!
		return ErrNotAuthorizedToModifyBalanceFromSvc
	}
	// assert correct command handling by the domain
	if ok := mbfs.Handle(ctx, a, &h.svc); !ok {
		var domErr error
		for i, e := range mbfs.Errors() {
			if i == 0 {
				domErr = e
			} else {
				domErr = errwrap.Wrap(domErr, e)
			}
		}
		return ErrModifyBalanceFromSvcFailedInDomain
	}
	// save domain facts to storage
	saveErr := h.rw.SaveFacts(ctx, target, app.OffersFactKeeper(&mbfs))
	if saveErr != nil {
		return errwrap.Wrap(ErrModifyBalanceFromSvcSavingFailed, saveErr)
	}
	return nil
}

// compile time assertions
var (
	_ app.RequiresCommandHandler = (*domain.ModifyBalanceFromSvc)(nil)
	_ app.RequiresErrorKeeper    = (*domain.ModifyBalanceFromSvc)(nil)
	_ app.OffersFactKeeper       = (*domain.ModifyBalanceFromSvc)(nil)
)
