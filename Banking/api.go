package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// ========== Utility ==========

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// ========== Error Handling Wrapper ==========

type apiFunc func(w http.ResponseWriter, r *http.Request) error

type ApiError struct {
	Error string `json:"error"`
}

func makeHandler(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

// ========== API Server ==========

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{listenAddr: listenAddr, store: store}
}

func (s *APIServer) Run() error {
	r := mux.NewRouter()

	r.HandleFunc("/account/transfer", makeHandler(s.handleTransferAccount)).Methods("POST")
	r.HandleFunc("/account/{id:[0-9]+}", makeHandler(s.handleGetAccountbyID)).Methods("GET")
	r.HandleFunc("/account/{id:[0-9]+}", makeHandler(s.handleDeleteAccount)).Methods("DELETE")
	r.HandleFunc("/account", makeHandler(s.handleCreateAccount)).Methods("POST")
	r.HandleFunc("/account", makeHandler(s.handleAccount)).Methods("GET")

	log.Println("JSON API server running on port", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, r)
}

// ========== Handlers ==========

func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	account := NewAccount(1, "John", "Doe", "1234567890", 1000.0, time.Now())
	writeJSON(w, http.StatusOK, account)
	return nil
}

func (s *APIServer) handleGetAccountbyID(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid account ID: %s", idStr)
	}

	account, err := s.store.GetAccount(id)
	if err != nil {
		return fmt.Errorf("error getting account: %v", err)
	}

	writeJSON(w, http.StatusOK, account)
	return nil
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	var req CreatAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return err
	}

	account := NewAccount(0, req.FirstName, req.LastName, req.Number, req.Balance, time.Now())
	if err := s.store.CreateAccount(account); err != nil {
		return err
	}

	writeJSON(w, http.StatusCreated, account)
	return nil
}

func (s *APIServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return fmt.Errorf("invalid account ID: %s", idStr)
	}

	err = s.store.DeleteAccount(id)
	if err != nil {
		return fmt.Errorf("failed to delete account: %v", err)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Account deleted successfully"))
	return nil
}

func (s *APIServer) handleTransferAccount(w http.ResponseWriter, r *http.Request) error {
	var req TransferAccountRequest

	// Parse and validate JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return fmt.Errorf("invalid request body: %v", err)
	}
	req.CreatedAt = time.Now()

	if req.FromAccountID == req.ToAccountID {
		return fmt.Errorf("cannot transfer to the same account")
	}
	if req.Amount <= 0 {
		return fmt.Errorf("amount must be greater than 0")
	}

	// Fetch sender and receiver accounts
	fromAccount, err := s.store.GetAccount(req.FromAccountID)
	if err != nil {
		return fmt.Errorf("from_account not found: %v", err)
	}

	toAccount, err := s.store.GetAccount(req.ToAccountID)
	if err != nil {
		return fmt.Errorf("to_account not found: %v", err)
	}

	// Check balance
	if fromAccount.Balance < req.Amount {
		return fmt.Errorf("insufficient balance in from_account")
	}

	// Update balances
	fromAccount.Balance -= req.Amount
	toAccount.Balance += req.Amount

	// Save updates
	if err := s.store.UpdateAccount(fromAccount); err != nil {
		return fmt.Errorf("failed to update from_account: %v", err)
	}
	if err := s.store.UpdateAccount(toAccount); err != nil {
		return fmt.Errorf("failed to update to_account: %v", err)
	}

	// Optional: log the transfer somewhere

	// Send response
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"message":          "transfer successful",
		"from_account_id":  fromAccount.ID,
		"to_account_id":    toAccount.ID,
		"from_balance":     fromAccount.Balance,
		"to_balance":       toAccount.Balance,
		"transferred_at":   req.CreatedAt,
	})

	return nil
}

