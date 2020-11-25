package requires

// DistinguishableAssertable can be asserted to be distinguishable
// application requires to be able to assert that Distinguishable can actually be identified
type DistinguishableAssertable interface {
	// IsDistinguishable knows how to assert that a potential Distinguishable can be actually identified
	IsDistinguishable() bool
}
