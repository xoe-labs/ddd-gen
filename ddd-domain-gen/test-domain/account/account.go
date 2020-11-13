package account

import (
	"github.com/xoe-labs/go-generators/ddd-domain-gen/test-domain/holder"
)

//go:generate go run ../../main.go -t Account
type Account struct {
	uuid    string        `ddd:"required'field uuid is empty'"`
	holder  holder.Holder `ddd:"required'field holder is empty'"`
	address string
	balance int64         `ddd:"private"` // read via domain logic: don't generate default getter
	values  []int64       `ddd:"private"`
}
