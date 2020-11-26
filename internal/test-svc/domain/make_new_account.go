package domain

//go:generate go run ../../../main.go --config ../ddd-config.yaml domain -t MakeNewAccount
type MakeNewAccount struct {
	errors []error
	facts []interface{}
}

