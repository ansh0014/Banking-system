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
	// TODO: Implement this later
	return nil
}
