package command

//go:generate go run ../../../../main.go --config ../../ddd-config.yaml app command --fact-based -t Commands
type Commands struct {
	MakeNewAccount          MakeNewAccountHandlerWrapper          ``
	MakeNewAccountWithOutId MakeNewAccountWithOutIdHandlerWrapper `command:"topic,account"`
	DeleteAccount           DeleteAccountHandlerWrapper           ``
	BlockAccount            BlockAccountHandlerWrapper            ``
	ValidateHolder          BlockAccountHandlerWrapper            `command:"w/o policy"`
	IncreaseBalance         IncreaseBalanceHandlerWrapper         ``
	IncreaseBalanceFromSvc  IncreaseBalanceHandlerWrapper         `command:"topic,balance; adapters,svc:github.com/xoe-labs/ddd-gen/internal/test-svc/app/ifaces.Balancer"`
}
