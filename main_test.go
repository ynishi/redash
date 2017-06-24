package redash

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	origUrlEnv    = os.Getenv(redashUrlEnv)
	origApikeyEnv = os.Getenv(redashApikeyEnv)
	testUrl       string
	testApikey    = "abcdef"
	server        *httptest.Server
)

func setup() {

	mux := http.NewServeMux()
	mux.HandleFunc(
		"/api/test",
		func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/api/test" {
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"result": "ok"}`)
		})
	server = httptest.NewServer(mux)
	os.Setenv(redashUrlEnv, server.URL)
	os.Setenv(redashApikeyEnv, testApikey)

}

func TestMain(m *testing.M) {
	setup()
	defer os.Setenv(redashUrlEnv, origUrlEnv)
	defer os.Setenv(redashApikeyEnv, origApikeyEnv)
	defer server.Close()
	v := m.Run()
	os.Exit(v)
}
