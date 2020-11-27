package app

import (
	"github.com/xoe-labs/ddd-gen/internal/test-svc/app/authorizable"
	"github.com/xoe-labs/ddd-gen/internal/test-svc/app/distinguishable"
)

// compile time assertions
var (
	_ RequiresDistinguishableAsserter = (*distinguishable.Target)(nil)
	_ OffersDistinguishable           = (*distinguishable.Target)(nil)
	_ OffersAuthorizable              = (*authorizable.Actor)(nil)
)
