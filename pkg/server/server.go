package server

import (
	"github.com/blokje5/iam-server/pkg/storage/database/postgres"
	"context"
	"github.com/blokje5/iam-server/pkg/storage"
	"net/http"

	"github.com/gorilla/mux"
)

// Server represents an instance of the IAM-authorize server
type Server struct {
	Handler http.Handler
	
	router *mux.Router
	NamespaceServer

	storage *storage.Storage
}

// New returns a new instance of the Server
func New() *Server {
	s := Server{}
	return &s
}

// Init initializes the server
func (s *Server) Init(ctx context.Context) error{
	r := mux.NewRouter()

	nr := r.PathPrefix("/namespaces").Subrouter()
	s.NamespaceServer.Init(nr)
	
	s.router = r

	pdb := postgres.New(postgres.PostgresConfig{})
	db, err := pdb.Initialize(ctx)
	if err != nil {
		return err
	}
	storage := storage.New(db)
	s.storage = storage
	return nil
}