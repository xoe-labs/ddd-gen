package domain

import (
	"github.com/xoe-labs/ddd-gen/internal/test-svc/domain/holder"
)

//go:generate go run ../../../main.go --config ../../ddd-config.yaml domain entity -t Account
type Account struct {
	uuid        string                   `entity:"required,field uuid is empty;equal;stringer"`
	holder      holder.Holder            `entity:"required,field holder is empty;getter"`
	altHolders  []holder.Holder          `entity:"required,field alternative holders is empty;setter"`
	holderRoles map[holder.Holder]string `entity:"required,field holder role map is empty"`
	address     string
	balance     int64   `entity:"private"` // read via domain logic: don't generate default getter
	values      []int64 `entity:"private"`
}

