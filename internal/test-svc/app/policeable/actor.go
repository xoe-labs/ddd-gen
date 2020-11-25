package policeable

import (
	// "fmt"

	"github.com/xoe-labs/ddd-gen/internal/test-svc/app"

	// "github.com/satori/go.uuid"
)

// Actor represents a target entity    --- this message component probably should be autogenerated, eg. by protobuf
type Actor struct {
	user           string
	elevationToken string
}

func (a *Actor) User() string {
	return a.user
}

func (a *Actor) ElevationToken() string {
	return a.elevationToken
}

var (
	_ app.OffersPoliceable = (*Actor)(nil)
)