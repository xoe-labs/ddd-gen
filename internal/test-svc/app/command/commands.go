package command

//go:generate go run ../../../../main.go --config ../../ddd-config.yaml app command --fact-based -t Commands
type Commands struct {
	MakeNewAccount       MakeNewAccountHandlerWrapper     ``
	MakeNewAccountQuick  MakeNewAccountQuckHandlerWrapper `command:"topic,account"`
	ArchiveAccount       ArchiveAccountHandlerWrapper     ``
	BlockAccount         BlockAccountHandlerWrapper       ``
	ValidateHolder       BlockAccountHandlerWrapper       `command:"w/o policy"`
	ModifyBalance        ModifyBalanceHandlerWrapper      ``
	ModifyBalanceFromSvc ModifyBalanceHandlerWrapper      `command:"topic,balance; adapters,svc:github.com/xoe-labs/ddd-gen/internal/test-svc/app/ifaces.Balancer"`
}
