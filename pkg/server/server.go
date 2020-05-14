package server

import (
	"context"
	"github.com/blokje5/iam-server/pkg/storage"
	"github.com/blokje5/iam-server/pkg/storage/database/postgres"
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
	params  *Params
}

// New returns a new instance of the Server
func New(params *Params) *Server {
	s := Server{
		params: params,
	}
	return &s
}

// Init initializes the server
func (s *Server) Init(ctx context.Context) error {
	r := mux.NewRouter()
	s.router = r

	pgConfig := postgres.NewConfig().SetConnectionString(s.params.ConnectionString)

	pdb := postgres.New(pgConfig)
	db, err := pdb.Initialize(ctx)
	if err != nil {
		return err
	}
	storage := storage.New(db)
	s.storage = storage

	nr := r.PathPrefix("/namespaces").Subrouter()
	s.NamespaceServer.Init(nr, storage)
	return nil
}

// Run runs the server until it is stopped
func (s *Server) Run(ctx context.Context) error {
	if err := http.ListenAndServe(":8080", s.router); err != nil {
		return err
	}

	return nil
}