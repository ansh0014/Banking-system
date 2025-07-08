package main

import (
	"database/sql"
	"log"

	// "time"

	_ "github.com/lib/pq"
)

// ========== Storage Interface ==========
type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccount(int) (*Account, error)
	GetAccountsbyID(int) ([]*Account, error)
	CreateAccountTable() error
}

// ========== DB Struct ==========
type PostgresStore struct {
	db *sql.DB
}

// ========== Init DB ==========
func initDB() *sql.DB {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}
	db, err := sql.Open("postgres", config.PostgresURL)
	if err != nil {
		log.Fatal(err)
	}
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

// ========== Create Table ==========
func (s *PostgresStore) CreateAccountTable() error {
	query := `
	CREATE TABLE IF NOT EXISTS account (
		id SERIAL PRIMARY KEY,
		first_name VARCHAR(50) NOT NULL,
		last_name VARCHAR(50) NOT NULL,
		number VARCHAR(20) NOT NULL UNIQUE,
		balance DECIMAL(10, 2) NOT NULL DEFAULT 0.00,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`
	_, err := s.db.Exec(query)
	return err
}

// ========== Create ==========
func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `INSERT INTO account (first_name, last_name, number, balance, created_at) 
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := s.db.Exec(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt)
	return err
}

// ========== Delete ==========
func (s *PostgresStore) DeleteAccount(id int) error {
	query := `DELETE FROM account WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}

// ========== Update ==========
func (s *PostgresStore) UpdateAccount(acc *Account) error {
	query := `UPDATE account SET first_name = $1, last_name = $2, number = $3, balance = $4, created_at = $5 WHERE id = $6`
	_, err := s.db.Exec(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt, acc.ID)
	return err
}

// ========== Get by ID ==========
func (s *PostgresStore) GetAccount(id int) (*Account, error) {
	query := `SELECT id, first_name, last_name, number, balance, created_at FROM account WHERE id = $1`
	row := s.db.QueryRow(query, id)
	return ScanAccount(row)
}

// ========== Get All with Same ID ==========
func (s *PostgresStore) GetAccountsbyID(id int) ([]*Account, error) {
	query := `SELECT id, first_name, last_name, number, balance, created_at FROM account WHERE id = $1`
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

// ========== Scan Helpers ==========

func ScanAccount(row *sql.Row) (*Account, error) {
	var acc Account
	err := row.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Number, &acc.Balance, &acc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}
