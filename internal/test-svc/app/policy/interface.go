package policy

import (
	"context"
)

type Policer interface {
	Can(ctx context.Context, p Policeable, action string, data []byte) bool
}

type Policeable interface {
	User() string
	ElevationToken() string
}
