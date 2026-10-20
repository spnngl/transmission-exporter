package transmission

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

const sessionHeader = "X-Transmission-Session-Id"

// fakeTransmission answers 409 with a new session id until the request carries it.
func fakeTransmission(t *testing.T, status int) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(sessionHeader) != "current" {
			w.Header().Set(sessionHeader, "current")
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(status)
		_, _ = w.Write([]byte(`{"arguments":{"version":"4.0.6"},"result":"success"}`))
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestGetSessionRenewsExpiredSessionID(t *testing.T) {
	c := New(fakeTransmission(t, http.StatusOK).URL, nil)
	c.sessionID = "expired"

	s, err := c.GetSession()
	if err != nil {
		t.Fatal(err)
	}
	if s.Version != "4.0.6" {
		t.Errorf("version = %q, want 4.0.6", s.Version)
	}
}

func TestGetSessionUnauthorized(t *testing.T) {
	c := New(fakeTransmission(t, http.StatusUnauthorized).URL, nil)

	if _, err := c.GetSession(); err == nil {
		t.Fatal("want an error on 401")
	}
}
