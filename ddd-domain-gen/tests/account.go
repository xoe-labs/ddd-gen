package tests

//go:generate go run ../main.go -t Account
type Account struct {
	uuid    *string  `gen:"getter" ddd:"required'field uuid is missing'"`
	holder  *string  `gen:"getter" ddd:"required'field holder is missing'"`
	address *string  `gen:"getter"`
	balance *int64   `ddd:"private"` // read via domain logic: don't generate default getter
	values  *[]int64 `ddd:"private" gen:"getter"`
}
