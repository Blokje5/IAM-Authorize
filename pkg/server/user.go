package server

import (
	"encoding/json"
	"net/http"

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
