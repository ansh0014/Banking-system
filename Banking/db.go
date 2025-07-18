package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// ===== Interface for Account Operations =====
type Storage interface {
	CreateAccount(*Account) error
	DeleteAccount(int) error
	UpdateAccount(*Account) error
	GetAccount(int) (*Account, error)
	GetAccountsByID(int) ([]*Account, error)
	CreateAccountTable() error
	Userlogin(username, password string) (*User, error)
	CreateUser(*User) error
	CreateUsersTable() error // This is for table creation only
}

// / ===== DB Wrapper Struct =====
type PostgresStore struct {
	db *sql.DB
}

// ===== Initialize DB Connection =====
func initDB() *sql.DB {
	config, err := LoadConfig()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	db, err := sql.Open("postgres", config.PostgresURL)
	if err != nil {
		log.Fatal("Database open error:", err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal("Database ping error:", err)
	}

	return db
}

// ===== Create Account Table =====
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






// ===== Insert New Account =====
func (s *PostgresStore) CreateAccount(acc *Account) error {
	query := `
		INSERT INTO account (first_name, last_name, number, balance, created_at) 
		VALUES ($1, $2, $3, $4, $5)`
	if acc.CreatedAt.IsZero() {
		acc.CreatedAt = time.Now()
	}
	_, err := s.db.Exec(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt)
	return err
}

// ===== Delete Account by ID =====
func (s *PostgresStore) DeleteAccount(id int) error {
	query := `DELETE FROM account WHERE id = $1`
	_, err := s.db.Exec(query, id)
	return err
}

// ===== Update Existing Account =====
func (s *PostgresStore) UpdateAccount(acc *Account) error {
	query := `
		UPDATE account 
		SET first_name = $1, last_name = $2, number = $3, balance = $4, created_at = $5 
		WHERE id = $6`
	_, err := s.db.Exec(query, acc.FirstName, acc.LastName, acc.Number, acc.Balance, acc.CreatedAt, acc.ID)
	return err
}

// ===== Get One Account by ID =====
func (s *PostgresStore) GetAccount(id int) (*Account, error) {
	query := `SELECT id, first_name, last_name, number, balance, created_at FROM account WHERE id = $1`
	row := s.db.QueryRow(query, id)
	return scanAccount(row)
}

// ===== Get All Accounts with Same ID (if needed) =====
func (s *PostgresStore) GetAccountsByID(id int) ([]*Account, error) {
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

func (s *PostgresStore) CreateUsersTable() error {
	query := `CREATE TABLE IF NOT EXISTS users (
        id SERIAL PRIMARY KEY,
        username VARCHAR(50) NOT NULL UNIQUE,
        password VARCHAR(100) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    )`
	_, err := s.db.Exec(query)
	return err
}

// ===== Get User by Username for Login =====
func (s *PostgresStore) Userlogin(username, password string) (*User, error) {
	query := `SELECT id, username, password FROM users WHERE username = $1`
	row := s.db.QueryRow(query, username)

	var user User
	err := row.Scan(&user.ID, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}

	// Verify password using bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return &user, nil
}

// ===== Insert New User =====
func (s *PostgresStore) CreateUser(user *User) error {
	query := `
        INSERT INTO users (username, password, created_at)
        VALUES ($1, $2, $3)
        RETURNING id`

	err := s.db.QueryRow(
		query,
		user.Username,
		user.Password,
		time.Now(),
	).Scan(&user.ID)

	if err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}
	return nil
}

// ===== Row Scan Helper =====
func scanAccount(row *sql.Row) (*Account, error) {
	var acc Account
	err := row.Scan(&acc.ID, &acc.FirstName, &acc.LastName, &acc.Number, &acc.Balance, &acc.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &acc, nil
}
