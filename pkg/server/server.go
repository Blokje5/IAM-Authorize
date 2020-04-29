package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

// Server represents an instance of the IAM-authorize server
type Server struct {
	Handler http.Handler
	
	router *mux.Router
	NamespaceServer
}

// New returns a new instance of the Server
func New() *Server {
	s := Server{}
	return &s
}

// Init initializes the server
func (s *Server) Init() {
	r := mux.NewRouter()

	nr := r.PathPrefix("/namespaces").Subrouter()
	s.NamespaceServer.Init(nr)
	
	s.router = r
}