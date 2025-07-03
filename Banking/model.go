package main
type Account struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Balance float64 `json:"balance"`

}
func NewAccount(id, name string, balance float64) *Account {
	return &Account{
		ID:      id,
		Name:    name,
		Balance: balance,
	}
}