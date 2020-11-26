package app

// OffersAuthorizable is an actor that can be policed
// application implements OffersAuthorizable and thereby offers policy adapter and external consumers a common language to reason about a authorizable actor
// TODO: implement OffersAuthorizable
type OffersAuthorizable interface {
	// TODO: adapt to your needs

	User() string
	ElevationToken() string
}
