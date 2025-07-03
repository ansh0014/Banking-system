package main

import "time"

type Account struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Number    string    `json:"number"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

func NewAccount(id int, firstName, lastName, number string, balance float64, createdAt time.Time) *Account {
	return &Account{
		ID:        id,
		FirstName: firstName,
		LastName:  lastName,
		Number:    number,
		Balance:   balance,
		CreatedAt: createdAt,
	}
}
