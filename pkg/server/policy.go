package server

import (
	"encoding/json"
	"net/http"

	"github.com/blokje5/iam-server/pkg/log"
	"github.com/blokje5/iam-server/pkg/server/middleware"
	"github.com/blokje5/iam-server/pkg/storage"
	"github.com/gorilla/mux"
)

type PolicyServer struct {
	Handler http.Handler
	router  *mux.Router
	storage *storage.Storage

	logger  *log.Logger
	decoder json.Decoder
}

// NewPolicyServer returns a new instance of the Policy Server
func NewPolicyServer() *PolicyServer {
	s := PolicyServer{
		logger: log.GetLogger(),
	}
	return &s
}

// Init initializes the server
func (s *PolicyServer) Init(r *mux.Router, middleware middleware.Middleware, storage *storage.Storage) {

	r.Handle("/", middleware(http.HandlerFunc(s.PostNamespaceHandler))).Methods("POST")

	s.router = r
	s.storage = storage
}

// PostNamespaceHandler handles Post requests on the policy resource
func (s *PolicyServer) PostNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)
	var policy *storage.Policy
	err := decoder.Decode(&policy)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	policy, err = s.storage.InsertPolicy(ctx, policy)
	if err != nil {
		if err == storage.ErrUniqueViolation {
			http.Error(w, NewConflictError("Conflict", "Uniqueness constraint violation").Error(), http.StatusConflict)
		} else {
			http.Error(w, NewInternalServerError("Internal server error", err.Error()).Error(), http.StatusInternalServerError)
		}
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(policy)
	if err != nil {
		http.Error(w, NewInternalServerError("Internal server error", err.Error()).Error(), http.StatusInternalServerError)
		return
	}
}