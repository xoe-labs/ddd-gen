package holder

import (
	"fmt"
	"time"
)

type HolderType struct{ s string }

func (p HolderType) String() string { return p.s }

var (
	Local  = HolderType{"local"}
	Remote = HolderType{"remote"}
)

//go:generate go run ../../../../main.go --config ../../ddd-config.yaml domain entity -t Holder -v validate
type Holder struct {
	name    string     `entity:"required,field name is empty;stringer"`
	altname string     ``
	bday    time.Time  ``
	hTyp    HolderType `entity:"required,filed folder type is empty"`
}

func (h Holder) validate() error {
	if h.altname == h.name {
		return fmt.Errorf("altname equals name")
	}
	return nil
}

// Apply applies facts to Holder
// implements application layer's entity interface.
func (h *Holder) Apply(fact interface{}) {
	// TODO: ipmlement
}
