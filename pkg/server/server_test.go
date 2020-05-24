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
	params.ConnectionString = connString
	params.MigrationPath = "../storage/database/postgres/migrations"

	server := New(params)
	if err := server.Init(ctx); err != nil {
		t.Fatalf("Could not start server: %v", err)
	}
	
	recorder := httptest.NewRecorder()

	return &fixture{
		server:   server,
		recorder: recorder,
		t:        t,
	}
}

func (f *fixture) executeRequestForHandler(handler http.Handler, req *http.Request, code int, resp string) {
	f.t.Helper()
	handler.ServeHTTP(f.recorder, req)
	if f.recorder.Code != code {
		f.t.Errorf("Expected status code %v, instead got: %v", code, f.recorder.Code)
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
