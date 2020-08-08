package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/blokje5/iam-server/pkg/log"
	"github.com/blokje5/iam-server/pkg/server/middleware"
	"github.com/blokje5/iam-server/pkg/storage"
	"github.com/gorilla/mux"
)

type UserServer struct {
	Handler http.Handler
	router  *mux.Router
	storage *storage.Storage

	logger *log.Logger
}

// NewUserServer returns a new instance of the User Server
func NewUserServer() *UserServer {
	s := UserServer{
		logger: log.GetLogger(),
	}
	return &s
}

func (s *UserServer) Init(r *mux.Router, middleware middleware.Middleware, storage *storage.Storage) {
	r.Handle("/", middleware(http.HandlerFunc(s.PostUserHandler))).Methods("POST")
	r.Handle("/{id:[0-9]+}", middleware(http.HandlerFunc(s.GetUserHandler))).Methods("GET")
	r.Handle("/{id:[0-9]+}", middleware(http.HandlerFunc(s.DeleteUserHandler))).Methods("DELETE")
	r.Handle("/{id:[0-9]+}/policies/{policy_id:[0-9]+}", middleware(http.HandlerFunc(s.PostUserPolicyHandler))).Methods("POST")
	r.Handle("/{id:[0-9]+}/policies/", middleware(http.HandlerFunc(s.GetUserPolicyHandler))).Methods("GET")

	s.router = r
	s.storage = storage
}

// PostUserHandler handles Post requests on the User resource
func (s *UserServer) PostUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)
	var user *storage.User
	err := decoder.Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err = s.storage.InsertUser(ctx, user)
	if err != nil {
		if err == storage.ErrUniqueViolation {
			http.Error(w, NewConflictError("Conflict", "Uniqueness constraint violation").Error(), http.StatusConflict)
		} else {
			http.Error(w, NewInternalServerError("Internal server error", err.Error()).Error(), http.StatusInternalServerError)
		}
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(user)
	if err != nil {
		http.Error(w, NewInternalServerError("Internal server error", err.Error()).Error(), http.StatusInternalServerError)
		return
	}
}

// GetUserHandler handles Get requests on the User resource by ID
func (s *UserServer) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	ID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.storage.GetUser(ctx, ID)
	if err != nil || user == nil {
		http.Error(w, NewNotFoundError("Not Found", fmt.Sprintf("User with ID: %v not found", ID)).Error(), http.StatusNotFound)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(user)
	if err != nil {
		http.Error(w, NewInternalServerError("Internal server error", err.Error()).Error(), http.StatusInternalServerError)
		return
	}
}

// DeleteUserHandler handles Get requests on the User resource
func (s *UserServer) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	ID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.storage.DeleteUser(ctx, ID)
	if err != nil {
		http.Error(w, NewInternalServerError("Internal server error", err.Error()).Error(), http.StatusInternalServerError)
	}

	w.WriteHeader(204)
}

func (s *UserServer) PostUserPolicyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	ID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	policyID, err := strconv.ParseInt(vars["policy_id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = s.storage.InsertPolicyForUser(ctx, ID, policyID)
	if err != nil {
		if err == storage.ErrForeignKeyViolation {
			http.Error(w, NewBadRequest("Bad Request", err.Error()).Error(), http.StatusBadRequest)
		} else {
			http.Error(w, NewInternalServerError("Internal server error", err.Error()).Error(), http.StatusInternalServerError)
		}
	}

	w.WriteHeader(200)
}

// GetUserPolicyHandler handles Get requests on the user policy resource by ID
func (s *UserServer) GetUserPolicyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	ID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, err := s.storage.GetUser(ctx, ID)
	if err != nil || user == nil {
		http.Error(w, NewNotFoundError("Not Found", fmt.Sprintf("User with ID: %v not found", ID)).Error(), http.StatusNotFound)
		return
	}

	policies, err := s.storage.GetPoliciesForUser(ctx, user)
	if err != nil {
		http.Error(w, NewInternalServerError("Internal server error", err.Error()).Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(policies)
	if err != nil {
		http.Error(w, NewInternalServerError("Internal server error", err.Error()).Error(), http.StatusInternalServerError)
		return
	}
}
