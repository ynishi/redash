package redash

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestDefaultClient(t *testing.T) {

	client := DefaultClient

	u, err := client.Url()
	if err != nil {
		t.Error(err)
	}
	if u.String() != server.URL {
		t.Fatalf("Url is bad. want: %q, have: %q", server.URL, u.String())
	}

	apikey, err := client.Apikey()
	if err != nil {
		t.Error(err)
	}
	if apikey != testApikey {
		t.Fatalf("Apikey is bad. want: %q, have: %q", testApikey, apikey)
	}

	if client.HTTPClient() != http.DefaultClient {
		t.Fatalf("HTTPClient is bad. want: %q, have: %q", http.DefaultClient, client.HTTPClient())
	}

	if !reflect.DeepEqual(client.DefaultOpts(), &defaultOpts) {
		t.Fatalf("DefaultOptions is bad. want: %q, have: %q", &defaultOpts, client.DefaultOpts())
	}
}

type Check struct {
	Result string `json:"result"`
}

func TestGet(t *testing.T) {

	params := map[string]string{"id": "abc"}
	resp, err := Get("api/test", params)
	print(err)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var have Check
	buf, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(buf, &have)
	if want := "ok"; have.Result != want {
		t.Fatalf("resp is bad. want: %q, have: %q", want, have)
	}
}
