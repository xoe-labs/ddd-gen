package account

import (
	"github.com/xoe-labs/ddd-gen/internal/test-svc/app/distinguishable"
	"github.com/xoe-labs/ddd-gen/internal/test-svc/domain/holder"
)

//go:generate go run ../../../../main.go --config ../../ddd-config.yaml domain entity -t Account
type Account struct {
	holder         holder.Holder                            `entity:"required,field holder is empty;getter"`
	altHolders     map[distinguishable.Target]holder.Holder `entity:"required,field alternative holders is empty;setter"`
	holderRoles    map[holder.Holder]string                 `entity:"required,field holder role map is empty"`
	address        string
	balance        int64   `entity:"private"` // read via domain logic: don't generate default getter
	movements      []int64 `entity:"private"`
	blocked        bool    `entity:"private"`
	blockReason    bool    `entity:"private"`
	unblockReason  bool    `entity:"private"`
	archived       bool    `entity:"private"`
	archivedReason bool    `entity:"private"`
	validHolder    bool    `entity:"private"`
	validatedBy    string  `entity:"private"`
}

// Apply applies facts to Account
// implements application layer's entity interface.
func (a *Account) Apply(fact interface{}) {
	// TODO: ipmlement
}

