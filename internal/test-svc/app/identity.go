package app

// RequiresDistinguishableAssertable can be asserted to be distinguishable
// application requires to be able to assert that OffersDistinguishable can actually be identified
type RequiresDistinguishableAssertable interface {
	// IsDistinguishable knows how to assert that a potential OffersDistinguishable can be actually identified
	IsDistinguishable() bool
}
