package main

import (
	"encoding/json"
	"fmt"

	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"

	

	_ "github.com/lib/pq"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		println("Error encoding JSON response:", err)
	}
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error
type ApiError struct {
	Error string
}

func makeHandler(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})

		}
	}
}

type APIServer struct {
	listenAddr string
	store      Storage
}

func (s *APIServer) Run() error {
	r := mux.NewRouter()

	// Order from most specific to least specific
	r.HandleFunc("/account/transfer", makeHandler(s.handleTransferAccount)).Methods("POST")
	r.HandleFunc("/account/{id}", makeHandler(s.handleGetAccount)).Methods("GET")
	r.HandleFunc("/account/{id}", makeHandler(s.handleDeleteAccount)).Methods("DELETE")
	r.HandleFunc("/account", makeHandler(s.handleCreateAccount)).Methods("POST")
	r.HandleFunc("/account", makeHandler(s.handleAccount)).Methods("GET")

	log.Println("JSON API server running on port ", s.listenAddr)
	return http.ListenAndServe(s.listenAddr, r)
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

// Handler examples
func (s *APIServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	account := NewAccount(1, "John", "Doe", "1234567890", 1000.0, time.Now())
	writeJSON(w, http.StatusOK, account)
	return nil
}

func (s *APIServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	account, err := s.store.GetAccount(1) // Example ID, replace with actual logic
	if err != nil {
		return err
	}
	writeJSON(w, http.StatusOK, account)
	return nil
}

func (s *APIServer) handleGetAccountbyID(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	fmt.Print("id:", id)

	writeJSON(w, http.StatusOK, &Account{})
	return nil
}

func (s *APIServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	var req CreatAccountRequest
	w.Header().Set("Content-Type", "application/json")
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
	if err := r.ParseForm(); err != nil {
		return err
	}
	return nil
}

func (s *APIServer) handleTransferAccount(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseForm(); err != nil {
		return err
	}
	return nil
}
