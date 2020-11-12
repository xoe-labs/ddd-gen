package holder

import (
	"time"
)

type HolderType struct{ s string }

func (p HolderType) String() string { return p.s }

var (
	Local  = HolderType{"local"}
	Remote = HolderType{"remote"}
)


//go:generate go run ../../main.go -t Holder
type Holder struct {
	uuid *string     `gen:"getter" ddd:"required'field uuid is missing'"`
	name *string     `gen:"getter" ddd:"required'field name is missing'"`
	bday *time.Time  `gen:"getter"`
	hTyp *HolderType `gen:"getter" ddd:"required'filed folder type is missing'"`
}
