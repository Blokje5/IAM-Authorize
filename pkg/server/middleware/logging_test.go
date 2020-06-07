package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/blokje5/iam-server/pkg/log"
)

func Test_loggingHandler_ServeHTTP(t *testing.T) {
	var b bytes.Buffer

	testLogger := log.NewWithFlags(&b, 0)
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	middleware := loggingHandler{
		h:      testHandler,
		logger: testLogger,
	}

	req := httptest.NewRequest("GET", "/namespaces/1", nil)
	req.RemoteAddr = "127.0.0.1"
	rr := httptest.NewRecorder()
	middleware.ServeHTTP(rr, req)

	if rr.Code != 200 {
		t.Fatalf("l.ServeHTTP(): Expected status code 200, instead got: %v", rr.Code)
	}

	got := b.String()
	want := "127.0.0.1 GET /namespaces/1\n"
	if got != want {
        t.Errorf("l.ServeHTTP() = %q, want %q", got, want)
    }
}
