package offers

import requires "github.com/xoe-labs/ddd-gen/internal/test-svc/app/ifaces/requires"

// Distinguishable can be identified
// application implements Distinguishable and thereby offers storage adapter and external consumers a common language to reason about identity
// TODO: implement Distinguishable
type Distinguishable interface {
	requires.DistinguishableAssertable
	// Identifier knows how to identify Distinguishable
	// TODO: adapt return type to your needs
	Identifier() string
}
