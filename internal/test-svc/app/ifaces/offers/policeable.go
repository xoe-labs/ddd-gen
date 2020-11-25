package offers

// Policeable is an actor that can be policed
// application implements Policeable and thereby offers policy adapter and external consumers a common language to reason about a policeable actor
// TODO: implement Policeable
type Policeable interface {
	// TODO: adapt to your needs

	User() string
	ElevationToken() string
}
