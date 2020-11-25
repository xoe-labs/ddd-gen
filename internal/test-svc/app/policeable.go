package app

// OffersPoliceable is an actor that can be policed
// application implements OffersPoliceable and thereby offers policy adapter and external consumers a common language to reason about a policeable actor
// TODO: implement OffersPoliceable
type OffersPoliceable interface {
	// TODO: adapt to your needs

	User() string
	ElevationToken() string
}
