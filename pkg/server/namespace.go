package server

import (
	"encoding/json"
	"net/http"

	"github.com/blokje5/iam-server/pkg/storage"
	"github.com/gorilla/mux"
)

type NamespaceServer struct {
	Handler http.Handler
	router  *mux.Router
	storage *storage.Storage

	decoder json.Decoder
}

// New returns a new instance of the Server
func NewNamespaceServer() *NamespaceServer {
	s := NamespaceServer{}
	return &s
}

// Init initializes the server
func (s *NamespaceServer) Init(r *mux.Router, storage *storage.Storage) {

	r.HandleFunc("/", s.ListNamespaceHandler).Methods("GET")
	r.HandleFunc("/", s.PostNamespaceHandler).Methods("POST")
	r.HandleFunc("/{id}", s.GetNamespaceHandler).Methods("GET")
	r.HandleFunc("/{id}", s.PutNamespaceHandler).Methods("PUT")
	r.HandleFunc("/{id}", s.DeleteNamespaceHandler).Methods("DELETE")

	s.router = r
	s.storage = storage
}

func (s *NamespaceServer) ListNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	namespaceList, err := s.storage.ListNamespaces(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(namespaceList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// PostNamespaceHandler handles Post requests on the namespace resource
func (s *NamespaceServer) PostNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	decoder := json.NewDecoder(r.Body)
	var namespace storage.Namespace
	err := decoder.Decode(&namespace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	namespace = s.storage.InsertNamespace(ctx, namespace)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(namespace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *NamespaceServer) GetNamespaceHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *NamespaceServer) PutNamespaceHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *NamespaceServer) DeleteNamespaceHandler(w http.ResponseWriter, r *http.Request) {

}
