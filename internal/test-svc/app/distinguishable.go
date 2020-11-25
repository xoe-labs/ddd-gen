package app

// OffersDistinguishable can be identified
// application implements OffersDistinguishable and thereby offers storage adapter and external consumers a common language to reason about identity
// TODO: implement OffersDistinguishable
type OffersDistinguishable interface {
	RequiresDistinguishableAssertable
	// Identifier knows how to identify OffersDistinguishable
	// TODO: adapt return type to your needs
	Identifier() string
}
