package distinguishable

import (
	"fmt"

	"github.com/xoe-labs/ddd-gen/internal/test-svc/app"

	"github.com/satori/go.uuid"
)

// Target represents a target entity    --- this message component probably should be autogenerated, eg. by protobuf
type Target struct {
	Continent string
	Zone      string
	Office    string
	Id        uuid.UUID
}

// Identifier implements the Distinguishable interface
func (t *Target) Identifier() string {
	return fmt.Sprintf("%s-%s-%s-%s", t.Continent, t.Zone, t.Office, t.Id)
}

// IsDistinguishable implements the DistinguishableAsserter interface used by the
// application layer to assert valid targets
func (t *Target) IsDistinguishable() bool {
	// in this example:
	//   - continent and zone are not required for a target to be distinguishable
	//   - for example, it might be used optionally for sharding or routing purposes
	return t.Office != "" && t.Id != uuid.Nil
}

var (
	_ app.RequiresDistinguishableAsserter = (*Target)(nil)
	_ app.OffersDistinguishable = (*Target)(nil)
)
