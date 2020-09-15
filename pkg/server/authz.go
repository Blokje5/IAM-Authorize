package server

import (
	"encoding/json"
	"net/http"

	"github.com/blokje5/iam-server/pkg/engine"
	"github.com/blokje5/iam-server/pkg/log"
	"github.com/blokje5/iam-server/pkg/server/middleware"

	"github.com/blokje5/iam-server/pkg/storage"
	"github.com/gorilla/mux"
)

type AuthzServer struct {
	Handler http.Handler
	router  *mux.Router
	storage *storage.Storage

	logger  *log.Logger
	decoder json.Decoder
	engine *engine.Engine
}

// NewAuthzServer returns a new instance of the Authz Server
func NewAuthzServer(enginge *engine.Engine) *AuthzServer {
	s := AuthzServer{
		logger: log.GetLogger(),
		engine: enginge,
	}
	return &s
}

// Init initializes the server
func (s *AuthzServer) Init(r *mux.Router, middleware middleware.Middleware, storage *storage.Storage, engine *engine.Engine) {

	r.Handle("/", middleware(http.HandlerFunc(s.PostAuthzHandler))).Methods("POST")

	s.router = r
	s.storage = storage
	s.engine = engine
}

// PostAuthzHandler handles Post requests on the authz resource, returning whether the input should be authorized
func (s *AuthzServer) PostAuthzHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)
	var input *engine.Input
	err := decoder.Decode(&input)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := s.engine.Query(ctx, *input)
	if err != nil {
		http.Error(w, NewInternalServerError("Internal server error", err.Error()).Error(), http.StatusInternalServerError)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(res)
	if err != nil {
		http.Error(w, NewInternalServerError("Internal server error", err.Error()).Error(), http.StatusInternalServerError)
		return
	}
}