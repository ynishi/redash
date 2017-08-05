package redash

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	mockApikey    = "mockApikey"
	mockClient    = mockClientData{}
)

type mockClientData struct {
	MockUrl string
}

func (mockClientData) Apikey() (apikey string, err error) {
	return mockApikey, nil
}

func (md mockClientData) Url() (u *url.URL, err error) {
	return url.Parse(md.MockUrl)
}

func (mockClientData) HTTPClient() *http.Client {
	return &http.Client{}
}

func (mockClientData) DefaultOpts() *Options {
	return defaultOpts()
}

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
