package server

import (
	"context"
	"strings"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
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

	resp :=  `{
		"id": 1,
		"name": "test"
	}`

	f.executeRequestForHandler(f.server.Handler, req, 200, resp)

	req, err = http.NewRequest("POST", "http://localhost:8080/namespaces/", strings.NewReader(body))
	if err != nil {
		t.Fatalf("Could not create request: %v", err)
	}

	resp =  `{
		"type":"http://localhost:8080/errors/conflict",
		"title":"Conflict",
		"detail":"Uniqueness constraint violation"
	}`

	f.executeRequestForHandler(f.server.Handler, req, 409, resp)
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
		f.t.Errorf("Expected status code %v, instead got: %v", code, f.recorder.Code)
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
