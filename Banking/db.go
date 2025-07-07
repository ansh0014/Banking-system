package main

import (
	// "log"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UploadAccount(*Account) error
	GetAccount(int) (*Account, error)
	GetAccountsbyID(int) ([]*Account, error)
}

// Storage implementation struct
type PostgresStore struct {
	db *sql.DB
}

func initDB() *sql.DB {
	db, err := sql.Open("postgres", "user=pqgotest dbname=pqgotest sslmode=verify-full")
	if err != nil {
		log.Fatal(err)
	}

	return db
}
func (s *PostgresStore) CreatAccountTable(acc *Account) error {
	query := `CREATE TABLE Account IF NOT EXISTS account (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		number VARCHAR(20) NOT NULL UNIQUE,
		balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := s.db.Exec(query)
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil

}

func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `INSERT INTO account (first_name, last_name, number, balance, created_at) VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt)
	return err
}
func (s *PostgresStore) DeleteAccount(id int) error {
	query := `DELETE FROM account WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}
func (s *PostgresStore) UploadAccount(acc *Account) error {
	query := `UPDATE account SET first_name = $1, last_name = $2, number = $3, balance = $4, created_at = $5 WHERE id = $6`
	_, err := s.db.Exec(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt, acc.ID)
	return err
}
func (s *PostgresStore) GetAccount(id int) (*Account, error) {
	query := `SELECT * FROM account WHERE id = $1`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return nil, nil
}
func (s *PostgresStore) GetAccountsbyID(id int) ([]*Account, error) {
	query := `SELECT * FROM account WHERE id = $1`
	rows, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var accounts []*Account
	for rows.Next() {
		var acc Account
		if err := rows.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Number, &acc.Balance, &acc.CreatedAt); err != nil {
			return nil, err
		}
		accounts = append(accounts, &acc)
	}

	return accounts, nil
}
func ScanAccount(rows *sql.Rows) (*Account, error) {
	var acc Account
	if err := rows.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Number, &acc.Balance, &acc.CreatedAt); err != nil {
		return nil, err
	}
	return &acc, nil
}

