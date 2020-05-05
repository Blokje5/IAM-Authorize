package server

import (
	"github.com/blokje5/iam-server/pkg/storage/database/postgres"
	"context"
	"github.com/blokje5/iam-server/pkg/storage"
	"net/http"

	"github.com/gorilla/mux"
)


// Params represents the configuration of the server
type Params struct {
	ConnectionString string
}

// NewParams returns a pointer to a new instance of the params struct
func NewParams() *Params {
	return &Params{}
}

// Server represents an instance of the IAM-authorize server
type Server struct {
	Handler http.Handler
	
	router *mux.Router
	NamespaceServer

	storage *storage.Storage
	params *Params
}

// New returns a new instance of the Server
func New(params *Params) *Server {
	s := Server{
		params: params,
	}
	return &s
}

// Init initializes the server
func (s *Server) Init(ctx context.Context) error{
	r := mux.NewRouter()

	nr := r.PathPrefix("/namespaces").Subrouter()
	s.NamespaceServer.Init(nr)
	
	s.router = r

	pdb := postgres.New(postgres.PostgresConfig{ConnectionString: s.params.ConnectionString})
	db, err := pdb.Initialize(ctx)
	if err != nil {
		return err
	}
	storage := storage.New(db)
	s.storage = storage
	return nil
}