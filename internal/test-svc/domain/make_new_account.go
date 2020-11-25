package domain

import "github.com/xoe-labs/ddd-gen/internal/test-svc/domain/Account"

//go:generate go run ../../../main.go --config ../ddd-config.yaml domain -t MakeNewAccount
type MakeNewAccount struct {
	errors []error
	facts []interface{}
}


// Handle handles MakeNewAccount in the domain
// returns true for success or false for failure.
// record errors with raise(err).
// implements application layer's CommandHandler interface.
func (mna *MakeNewAccount) Handle(a *account.Account) bool {
	return true
}
