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
	MakeNewAccount          command.MakeNewAccountHandlerWrapper          ``
	MakeNewAccountWithOutId command.MakeNewAccountWithOutIdHandlerWrapper `command:"topic,account"`
	DeleteAccount           command.DeleteAccountHandlerWrapper           ``
	BlockAccount            command.BlockAccountHandlerWrapper            ``
	ValidateHolder          command.BlockAccountHandlerWrapper            `command:"w/o policy"`
	IncreaseBalance         command.IncreaseBalanceHandlerWrapper         ``
	IncreaseBalanceFromSvc  command.IncreaseBalanceHandlerWrapper         `command:"topic,balance; adapters,svc:github.com/xoe-labs/ddd-gen/internal/test-svc/app/balancer.Balancer"`
}

type Queries struct {
	// HourAvailability      query.HourAvailabilityHandler
}
