package redash

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	defaultHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Ok")
	})

	origUrlEnv    = os.Getenv(redashUrlEnv)
	origApikeyEnv = os.Getenv(redashApikeyEnv)
	testUrl       string
	testApikey    = "abcdef"
	server        *httptest.Server
)

func setup() {

	server = httptest.NewServer(defaultHandler)
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
