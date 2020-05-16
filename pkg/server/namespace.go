package server

import (
	"encoding/json"
	"net/http"
	"strconv"

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
	r.HandleFunc("/{id:[0-9]+}", s.GetNamespaceHandler).Methods("GET")
	r.HandleFunc("/{id:[0-9]+}", s.PutNamespaceHandler).Methods("PUT")
	r.HandleFunc("/{id:[0-9]+}", s.DeleteNamespaceHandler).Methods("DELETE")

	s.router = r
	s.storage = storage
}

// ListNamespaceHandler handles GET requests on the namespace resource
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
	var namespace *storage.Namespace
	err := decoder.Decode(&namespace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	namespace = s.storage.InsertNamespace(ctx, *namespace)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(namespace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// GetNamespaceHandler handles Get by ID requests on the namespace resource
func (s *NamespaceServer) GetNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	namespace := s.storage.GetNamespace(ctx, ID)
	if namespace == nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(namespace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// PutNamespaceHandler handles PUT by ID requests on the namespace resource
func (s *NamespaceServer) PutNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	namespace := s.storage.GetNamespace(ctx, ID)
	if namespace == nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var newNamespace *storage.Namespace
	err = decoder.Decode(&newNamespace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	namespace.Name = newNamespace.Name

	newNamespace = s.storage.UpdateNamespace(ctx, *namespace)

	encoder := json.NewEncoder(w)
	err = encoder.Encode(newNamespace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// DeleteNamespaceHandler handles DELETE by ID requests on the namespace resource
func (s *NamespaceServer) DeleteNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	ID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.storage.DeleteNamespace(ctx, ID)
	w.Write([]byte("Success"))
}
