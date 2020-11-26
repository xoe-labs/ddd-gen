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


//go:generate go run ../../../../main.go --config ../../ddd-config.yaml domain entity -t Holder -v validate
type Holder struct {
	uuid string     `entity:"required,field uuid is empty;equal,reflect"`
	name string     `entity:"required,field name is empty;stringer"`
	bday time.Time
	hTyp HolderType `entity:"required,filed folder type is empty"`
}

func (h Holder) validate() error {
	if h.uuid == h.name {
		return fmt.Errorf("uuid equals name")
	}
	return nil
}

// Apply applies facts to Holder
// implements application layer's entity interface.
func (h *Holder) Apply(fact interface{}) {
	// TODO: ipmlement
}