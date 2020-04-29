package server

import (
	"net/http"

	"github.com/gorilla/mux"
)

type NamespaceServer struct {
	Handler http.Handler
	router *mux.Router
}

// New returns a new instance of the Server
func NewNamespaceServer() *NamespaceServer {
	s := NamespaceServer{}
	return &s
}

// Init initializes the server
func (s *NamespaceServer) Init(r *mux.Router) {

	r.HandleFunc("/namespaces", s.ListNamespaceHandler).Methods("GET")
	r.HandleFunc("/namespaces", s.PostNamespaceHandler).Methods("POST")
	r.HandleFunc("/namespaces/{id}", s.GetNamespaceHandler).Methods("GET")
	r.HandleFunc("/namespaces/{id}", s.PutNamespaceHandler).Methods("PUT")
	r.HandleFunc("/namespaces/{id}", s.DeleteNamespaceHandler).Methods("DELETE")

	s.router = r
}

func (s *NamespaceServer) ListNamespaceHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *NamespaceServer) PostNamespaceHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *NamespaceServer) GetNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	
}

func (s *NamespaceServer) PutNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	
}

func (s *NamespaceServer) DeleteNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	
}