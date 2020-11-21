package app

import (
	"github.com/xoe-labs/ddd-gen/internal/test-svc/app/command"
	// "github.com/xoe-labs/ddd-gen/internal/test-svc/app/query"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

//go:generate go run ../../../main.go --config ../ddd-config.yaml app command -t Commands
type Commands struct {
	MakeNewAccount          command.MakeNewAccountHandler          `command:"new"`
	MakeNewAccountWithOutId command.MakeNewAccountWithOutIdHandler `command:"new,non-identifiable; topic,account"`
	DeleteAccount           command.DeleteAccountHandler           `command:"del"`
	BlockAccount            command.BlockAccountHandler            ``
	ValidateHolder          command.BlockAccountHandler            `command:"w/o policy"`
	IncreaseBalance         command.IncreaseBalanceHandler         ``
	IncreaseBalanceFromSvc  command.IncreaseBalanceHandler         `command:"topic,balance; adapters,svc:github.com/xoe-labs/ddd-gen/internal/test-svc/app/balancer.Balancer"`
}

type Queries struct {
	// HourAvailability      query.HourAvailabilityHandler
}
