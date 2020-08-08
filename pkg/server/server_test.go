package server

import (
	"context"
	"encoding/json"
	"github.com/blokje5/iam-server/pkg/storage"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNamespaceServer_PostNamespaceHandler(t *testing.T) {
	f := newFixture(t)
	body := `{
		"name": "test"
	}`
	req, err := http.NewRequest("POST", "http://localhost:8080/namespaces/", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp := `{
		"id": 1,
		"name": "test"
	}`

	f.executeRequestForHandler(f.server.Handler, req, 200, resp)

	req, err = http.NewRequest("POST", "http://localhost:8080/namespaces/", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp = `{
		"type":"http://localhost:8080/errors/conflict",
		"status": 409,
		"title":"Conflict",
		"detail":"Uniqueness constraint violation"
	}`

	f.executeRequestForHandler(f.server.Handler, req, 409, resp)
}

func TestNamespaceServer_GetNamespaceHandler(t *testing.T) {
	f := newFixture(t)
	req, err := http.NewRequest("GET", "http://localhost:8080/namespaces/1", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp := `{
		"type":"http://localhost:8080/errors/not-found",
		"status": 404,
		"title":"Not Found",
		"detail":"Namespace with ID: 1 not found"
	}`

	f.executeRequestForHandler(f.server.Handler, req, 404, resp)

	_, err = f.server.storage.InsertNamespace(context.Background(), storage.Namespace{Name: "test"})
	if err != nil {
		t.Fatalf("Could not insert namespace: %v", err)
	}

	req, err = http.NewRequest("GET", "http://localhost:8080/namespaces/1", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp = `{
		"id": 1,
		"name": "test"
	}`

	f.executeRequestForHandler(f.server.Handler, req, 200, resp)
}

func TestNamespaceServer_ListNamespaceHandler(t *testing.T) {
	f := newFixture(t)
	req, err := http.NewRequest("GET", "http://localhost:8080/namespaces/", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp := `{
		"namespaces": []
	}`

	f.executeRequestForHandler(f.server.Handler, req, 200, resp)

	_, err = f.server.storage.InsertNamespace(context.Background(), storage.Namespace{Name: "test"})
	if err != nil {
		t.Fatalf("Could not insert namespace: %v", err)
	}

	req, err = http.NewRequest("GET", "http://localhost:8080/namespaces/", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp = `{
		"namespaces": [{
			"id": 1,
			"name": "test"
		}]
	}`

	f.executeRequestForHandler(f.server.Handler, req, 200, resp)
}

func TestNamespaceServer_DeleteNamespaceHandler(t *testing.T) {
	f := newFixture(t)
	_, err := f.server.storage.InsertNamespace(context.Background(), storage.Namespace{Name: "test"})
	if err != nil {
		t.Fatalf("Could not insert namespace: %v", err)
	}

	req, err := http.NewRequest("DELETE", "http://localhost:8080/namespaces/1", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp := ``

	f.executeRequestForHandler(f.server.Handler, req, 204, resp)
}

func TestNamespaceServer_PutNamespaceHandler(t *testing.T) {
	f := newFixture(t)
	_, err := f.server.storage.InsertNamespace(context.Background(), storage.Namespace{Name: "test"})
	if err != nil {
		t.Fatalf("Could not insert namespace: %v", err)
	}

	body := `{
		"name": "test2"
	}`

	req, err := http.NewRequest("PUT", "http://localhost:8080/namespaces/1", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp := `{
		"id": 1,
		"name": "test2"
	}`

	f.executeRequestForHandler(f.server.Handler, req, 200, resp)
}

func TestPolicyServer_GetPolicyHandler(t *testing.T) {
	f := newFixture(t)

	req, err := http.NewRequest("GET", "http://localhost:8080/policies/1", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp := `{
		"type":"http://localhost:8080/errors/not-found",
		"status": 404,
		"title":"Not Found",
		"detail":"Policy with ID: 1 not found"
	}`
	f.executeRequestForHandler(f.server.Handler, req, 404, resp)

	_, err = f.server.storage.InsertPolicy(context.Background(), storage.NewPolicy([]storage.Statement{storage.NewStatement(storage.Allow, []string{"*"}, []string{"iam:CreatePolicy"})}))
	if err != nil {
		t.Fatalf("Could not insert policy: %v", err)
	}

	req, err = http.NewRequest("GET", "http://localhost:8080/policies/1", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp = `{
		"id": 1,
		"version": "v1",
		"statements": [
		  {
			"Actions": [
			  "iam:CreatePolicy"
			],
			"Effect": "allow",
			"Resources": [
			  "*"
			]
		  }
		]
	  }`

	f.executeRequestForHandler(f.server.Handler, req, 200, resp)
}

func TestPolicyServer_PostNamespaceHandler(t *testing.T) {
	f := newFixture(t)
	body := `{
		"version": "v1",
		"statements": [
			{
				"effect": "allow",
				"actions": [
					"iam:CreatePolicy"
				],
				"resources": [
					"*"
				]
			}
		]
	}`
	req, err := http.NewRequest("POST", "http://localhost:8080/policies/", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp := `{
		"id": 1,
		"version": "v1",
		"statements": [
		  {
			"Actions": [
			  "iam:CreatePolicy"
			],
			"Effect": "allow",
			"Resources": [
			  "*"
			]
		  }
		]
	  }`

	f.executeRequestForHandler(f.server.Handler, req, 200, resp)
}

func TestPolicyServer_ListPolicyHandler(t *testing.T) {
	f := newFixture(t)
	req, err := http.NewRequest("GET", "http://localhost:8080/policies/", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp := `[]`

	f.executeRequestForHandler(f.server.Handler, req, 200, resp)

	_, err = f.server.storage.InsertPolicy(context.Background(), storage.NewPolicy([]storage.Statement{storage.NewStatement(storage.Allow, []string{"*"}, []string{"iam:CreatePolicy"})}))
	if err != nil {
		t.Fatalf("Could not insert policy: %v", err)
	}

	req, err = http.NewRequest("GET", "http://localhost:8080/policies/", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp = `[
		{
			"id": 1,
			"version": "v1",
			"statements": [
				{
				"Actions": [
					"iam:CreatePolicy"
				],
				"Effect": "allow",
				"Resources": [
					"*"
				]
				}
			]
		}
	]`

	f.executeRequestForHandler(f.server.Handler, req, 200, resp)
}

func TestPolicyServer_DeletePolicyHandler(t *testing.T) {
	f := newFixture(t)
	_, err := f.server.storage.InsertPolicy(context.Background(), storage.NewPolicy([]storage.Statement{storage.NewStatement(storage.Allow, []string{"*"}, []string{"iam:CreatePolicy"})}))
	if err != nil {
		t.Fatalf("Could not insert policy: %v", err)
	}

	req, err := http.NewRequest("DELETE", "http://localhost:8080/policies/1", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	f.executeRequestForHandler(f.server.Handler, req, 204, "")
}

func TestUserServer_PostUserHandler(t *testing.T) {
	f := newFixture(t)
	body := `{
		"name": "test"
	}`
	req, err := http.NewRequest("POST", "http://localhost:8080/users/", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp := `{
		"id": 1,
		"name": "test"
	}`

	f.executeRequestForHandler(f.server.Handler, req, 200, resp)

	req, err = http.NewRequest("POST", "http://localhost:8080/users/", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp = `{
		"type":"http://localhost:8080/errors/conflict",
		"status": 409,
		"title":"Conflict",
		"detail":"Uniqueness constraint violation"
	}`

	f.executeRequestForHandler(f.server.Handler, req, 409, resp)
}

func TestUsereServer_GetUsereHandler(t *testing.T) {
	f := newFixture(t)
	req, err := http.NewRequest("GET", "http://localhost:8080/users/1", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp := `{
		"type":"http://localhost:8080/errors/not-found",
		"status": 404,
		"title":"Not Found",
		"detail":"User with ID: 1 not found"
	}`

	f.executeRequestForHandler(f.server.Handler, req, 404, resp)

	_, err = f.server.storage.InsertUser(context.Background(), &storage.User{Name: "test"})
	if err != nil {
		t.Fatalf("Could not insert user: %v", err)
	}

	req, err = http.NewRequest("GET", "http://localhost:8080/users/1", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp = `{
		"id": 1,
		"name": "test"
	}`

	f.executeRequestForHandler(f.server.Handler, req, 200, resp)
}

func TestUserServer_DeleteUserHandler(t *testing.T) {
	f := newFixture(t)
	_, err := f.server.storage.InsertUser(context.Background(), &storage.User{Name: "test"})
	if err != nil {
		t.Fatalf("Could not insert user: %v", err)
	}

	req, err := http.NewRequest("DELETE", "http://localhost:8080/users/1", nil)
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	f.executeRequestForHandler(f.server.Handler, req, 204, "")
}


type fixture struct {
	server   *Server
	recorder *httptest.ResponseRecorder
	t        *testing.T
}

func newFixture(t *testing.T) *fixture {
	t.Helper()
	ctx := context.Background()

	params := NewParams()
	connString := os.Getenv("TEST_CONNECTIONSTRING")
	connString = "postgresql://postgres@127.0.0.1:5432/iam_test?sslmode=disable&password=local"
	params.ConnectionString = connString
	params.MigrationPath = "../storage/database/postgres/migrations"

	server := New(params)
	if err := server.Init(ctx); err != nil {
		t.Fatalf("Could not start server: %v", err)
	}

	recorder := httptest.NewRecorder()

	f := &fixture{
		server:   server,
		recorder: recorder,
		t:        t,
	}

	// Register test cleanup function
	t.Cleanup(f.cleanStorage)

	return f
}

func (f *fixture) reset() {
	f.recorder = httptest.NewRecorder()
}

func (f *fixture) cleanStorage() {
	ctx := context.Background()
	err := f.server.storage.Clean(ctx)
	if err != nil {
		f.t.Errorf("Failed to clean fixture: %v", err)
	}
}

func (f *fixture) executeRequestForHandler(handler http.Handler, req *http.Request, code int, resp string) {
	f.t.Helper()
	f.reset()
	handler.ServeHTTP(f.recorder, req)
	if f.recorder.Code != code {
		f.t.Errorf("Expected status code %v, instead got: %v. body: %s", code, f.recorder.Code, f.recorder.Body.String())
		return
	}
	if resp != "" {
		var result interface{}
		if err := json.Unmarshal(f.recorder.Body.Bytes(), &result); err != nil {
			f.t.Errorf("Expected JSON response from %v %v but got: %v", req.Method, req.URL, f.recorder)
		}
		var expected interface{}
		if err := json.Unmarshal([]byte(resp), &expected); err != nil {
			f.t.Fatalf("Unexpected error in expected response: %v", err)
		}
		if !reflect.DeepEqual(result, expected) {
			a, err := json.MarshalIndent(expected, "", "  ")
			if err != nil {
				f.t.Fatalf("Unexpected error in Marshal Indentiation: %v", err)
			}
			b, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				f.t.Fatalf("Unexpected error in Marshal Indentiation: %v", err)
			}

			f.t.Errorf("Expected JSON response from %v %v to equal:\n\n%s\n\nGot:\n\n%s", req.Method, req.URL, a, b)
		}
	}
}
