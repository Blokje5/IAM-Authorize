package server

import (
	"encoding/json"
	"fmt"
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
		return
	}

	namespace, err = s.storage.InsertNamespace(ctx, *namespace)
	if err != nil {
		if err == storage.ErrUniqueViolation {
			http.Error(w, NewConflictError("Conflict", "Uniqueness constraint violation").Error(), http.StatusConflict)
		} else {
			http.Error(w, NewInternalServerError("Internal server error", err.Error()).Error(), http.StatusInternalServerError)
		}
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(namespace)
	if err != nil {
		http.Error(w, NewInternalServerError("Internal server error", err.Error()).Error(), http.StatusInternalServerError)
		return
	}
}

// GetNamespaceHandler handles Get by ID requests on the namespace resource
func (s *NamespaceServer) GetNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	ID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	namespace, err := s.storage.GetNamespace(ctx, ID)
	if err != nil || namespace == nil {
		http.Error(w, NewNotFoundError("Not Found", fmt.Sprintf("Namespace with ID: %v not found", ID)).Error(), http.StatusNotFound)
		return
	}

	encoder := json.NewEncoder(w)
	err = encoder.Encode(namespace)
	if err != nil {
		http.Error(w, NewInternalServerError("Internal server error", err.Error()).Error(), http.StatusInternalServerError)
	}
}

// PutNamespaceHandler handles PUT by ID requests on the namespace resource
func (s *NamespaceServer) PutNamespaceHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	vars := mux.Vars(r)
	ID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	namespace, err := s.storage.GetNamespace(ctx, ID)
	if err != nil || namespace == nil {
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
	ID, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.storage.DeleteNamespace(ctx, ID)
	w.Write([]byte("Success"))
}
