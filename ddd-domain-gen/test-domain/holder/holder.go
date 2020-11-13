package holder

import (
	"time"
	"fmt"
)

type HolderType struct{ s string }

func (p HolderType) String() string { return p.s }

var (
	Local  = HolderType{"local"}
	Remote = HolderType{"remote"}
)


//go:generate go run ../../main.go -t Holder -v validate
type Holder struct {
	uuid string     `ddd:"required'field uuid is empty'"`
	name string     `ddd:"required'field name is empty'"`
	bday time.Time
	hTyp HolderType `ddd:"required'filed folder type is empty'"`
}

func (h Holder) validate() error {
	if h.uuid == h.name {
		return fmt.Errorf("uuid euqals name")
	}
	return nil
}
