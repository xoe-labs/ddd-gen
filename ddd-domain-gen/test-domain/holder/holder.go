package holder

import (
	"time"
)

//go:generate go run ../../main.go -t Holder
type Holder struct {
	uuid *string    `gen:"getter" ddd:"required'field uuid is missing'"`
	name *string    `gen:"getter" ddd:"required'field name is missing'"`
	bday *time.Time `gen:"getter"`
}
