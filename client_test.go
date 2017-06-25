package redash

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
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

	beforeUrlEnv := os.Getenv(redashUrlEnv)
	defer os.Setenv(redashUrlEnv, beforeUrlEnv)

	mux := http.NewServeMux()
	mux.HandleFunc(
		"/api/gettest",
		func(w http.ResponseWriter, r *http.Request) {
			if auth := r.Header.Get("Authorization"); !strings.Contains(auth, testApikey) {
				http.Error(w, fmt.Sprintf("Invalid Apikey %s", auth), http.StatusForbidden)
				return
			}
			if meth := r.Method; meth != http.MethodGet {
				http.Error(w, fmt.Sprintf("Invalid Method %s", meth), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"result": "ok"}`)
		})

	tgs := httptest.NewServer(mux)

	os.Setenv(redashUrlEnv, tgs.URL)

	params := map[string]string{"id": "abc"}
	resp, err := Get("api/gettest", params)
	if err != nil {
		t.Fatalf("Failed recieve Response from Get. %v", err)
	}
	defer resp.Body.Close()
	var have Check
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed ReadAll from buf, %v", err)
	}
	err = json.Unmarshal(buf, &have)
	if err != nil {
		t.Fatalf("Failed Decode Json, %v", err)
	}
	if resp.StatusCode == http.StatusForbidden {
		t.Fatalf("Apikey did not match, apikey: %s", resp.Body)
	}
	if resp.StatusCode == http.StatusBadRequest {
		t.Fatalf("HTTPMethod was not GET, method: %s", resp.Body)
	}
	if want := "ok"; have.Result != want {
		t.Fatalf("resp is bad. want: %q, have: %q", want, have)
	}
}

func TestPost(t *testing.T) {
	beforeUrlEnv := os.Getenv(redashUrlEnv)
	defer os.Setenv(redashUrlEnv, beforeUrlEnv)

	mux := http.NewServeMux()
	mux.HandleFunc(
		"/api/posttest",
		func(w http.ResponseWriter, r *http.Request) {
			if auth := r.Header.Get("Authorization"); !strings.Contains(auth, testApikey) {
				http.Error(w, fmt.Sprintf("Invalid Apikey %s", auth), http.StatusForbidden)
				return
			}
			if meth := r.Method; meth != http.MethodPost {
				http.Error(w, fmt.Sprintf("Invalid Method %s", meth), http.StatusBadRequest)
				return
			}
			buf, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Fatalf("Failed ReadAll have: %q err: %v", r.Body, err)
			}
			var postJson interface{}
			if err := json.Unmarshal(buf, &postJson); err != nil {
				t.Fatalf("Failed Json Decode have: %q err: %v", postJson, err)
			}
			if len(postJson.(interface{}).(map[string]interface{})["name"].(string)) == 0 {
				t.Fatalf("Length of postJson was 0 have: %q err: %v", postJson, err)
			}
			if postJson.(interface{}).(map[string]interface{})["name"].(string) != "abc" {
				http.Error(w, fmt.Sprintf("Failed send data have: %q %q", r.Body, postJson), http.StatusInternalServerError)
			}
			if !postJson.(interface{}).(map[string]interface{})["options"].(map[string]interface{})["opt1"].(bool) {
				http.Error(w, fmt.Sprintf("Failed send child data have: %q %q", r.Body, postJson), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"result": "ok"}`)
		})

	tgs := httptest.NewServer(mux)

	os.Setenv(redashUrlEnv, tgs.URL)

	reqMap := map[string]string{"id": "123", "name": "abc"}
	reqOpt := map[string]string{"opt1": "true", "opt2": "false"}
	resp, err := Post("api/posttest", reqMap, reqOpt)
	if err != nil {
		t.Fatalf("Failed recieve Response from Post. %v", err)
	}
	defer resp.Body.Close()
	var have Check
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed ReadAll from buf, %v", err)
	}
	if resp.StatusCode == http.StatusForbidden {
		t.Fatalf("Apikey did not match, apikey: %s", resp.Body)
	}
	if resp.StatusCode == http.StatusBadRequest {
		t.Fatalf("HTTPMethod was not POST, method: %s", resp.Body)
	}
	if resp.StatusCode == http.StatusInternalServerError {
		t.Fatalf("Cannot send body, BODY: %s", resp.Body)
	}
	err = json.Unmarshal(buf, &have)
	if err != nil {
		t.Fatalf("Failed Decode Json, %v", err)
	}
	if want := "ok"; have.Result != want {
		t.Fatalf("resp is bad. want: %q, have: %q", want, have)
	}
}

func TestDelete(t *testing.T) {}

func TestDo(t *testing.T) {}

func TestRequest(t *testing.T) {}

func TestOriginalClient(t *testing.T) {}

func TestGetInter(t *testing.T) {}

func TestPostInter(t *testing.T) {}

func TestDeleteInter(t *testing.T) {}

func TestDoInter(t *testing.T) {}

func TestRequestInter(t *testing.T) {}
