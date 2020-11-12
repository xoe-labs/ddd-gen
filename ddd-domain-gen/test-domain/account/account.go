package account

import (
	"github.com/xoe-labs/go-generators/ddd-domain-gen/test-domain/holder"
)

//go:generate go run ../../main.go -t Account
type Account struct {
	uuid    *string        `gen:"getter" ddd:"required'field uuid is missing'"`
	holder  *holder.Holder `gen:"getter" ddd:"required'field holder is missing'"`
	address *string        `gen:"getter"`
	balance *int64         `ddd:"private"` // read via domain logic: don't generate default getter
	values  *[]int64       `ddd:"private" gen:"getter"`
}
