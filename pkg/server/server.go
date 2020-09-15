package server

import (
	"context"
	"net/http"

	"github.com/blokje5/iam-server/pkg/log"
	"github.com/blokje5/iam-server/pkg/server/middleware"
	"github.com/blokje5/iam-server/pkg/storage"
	"github.com/blokje5/iam-server/pkg/storage/database/postgres"
	"github.com/blokje5/iam-server/pkg/engine"

	"github.com/gorilla/mux"
)

// Params represents the configuration of the server
type Params struct {
	ConnectionString string
	MigrationPath    string
	PolicyPath string
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
	PolicyServer
	UserServer

	logger  *log.Logger
	storage *storage.Storage
	params  *Params
	engine *engine.Engine
}

// New returns a new instance of the Server
func New(params *Params) *Server {
	s := Server{
		params: params,
		logger: log.GetLogger(),
	}
	return &s
}

// Init initializes the server
func (s *Server) Init(ctx context.Context) error {
	s.logger.Info("Initializing server")
	r := mux.NewRouter()
	s.router = r

	pgConfig := postgres.NewConfig().SetConnectionString(s.params.ConnectionString).SetMigrationPath(s.params.MigrationPath)
	s.logger.Debugf("Connecting to database with connection string: %s", s.params.ConnectionString)
	pdb := postgres.New(pgConfig)

	s.logger.Debug("Starting database initialization")
	db, err := pdb.Initialize(ctx)
	s.logger.Debug("Completed database initialization")

	if err != nil {
		return err
	}
	storage := storage.New(db, pdb)
	s.storage = storage

	s.logger.Debug("Starting engine initialization")
	engine, err := engine.Load([]string{s.params.PolicyPath})
	if err != nil {
		return err
	}
	s.engine = engine
	s.logger.Debug("Completed engine initialization")

	s.logger.Debug("Initializing routers")
	nr := r.PathPrefix("/namespaces").Subrouter()
	middleware := middleware.NewLoggingMiddleware(s.logger)
	s.NamespaceServer.Init(nr, middleware, storage)

	pr := r.PathPrefix("/policies").Subrouter()
	s.PolicyServer.Init(pr, middleware, storage)

	ur := r.PathPrefix("/users").Subrouter()
	s.UserServer.Init(ur, middleware, storage)
	s.logger.Debug("Completed Initializing routers")

	s.Handler = r
	return nil
}

// Run runs the server until it is stopped
func (s *Server) Run(ctx context.Context) error {
	if err := http.ListenAndServe(":8080", s.router); err != nil {
		return err
	}

	return nil
}
