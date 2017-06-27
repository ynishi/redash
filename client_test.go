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

	"bytes"
	"net/url"
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
			if id := r.URL.Query().Get("id"); id != "abc" {
				http.Error(w, fmt.Sprintf("Invalid Parameter %s", id), http.StatusInternalServerError)
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

type postOpt struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Opts Opt    `json:"opts"`
}
type Opt struct {
	Opt1 bool   `json:"opt1"`
	Opt2 string `json:"opt2"`
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
			var postJson postOpt
			if err := json.Unmarshal(buf, &postJson); err != nil {
				t.Fatalf("Failed Json Decode have: %q err: %v", postJson, err)
			}
			if postJson == *new(postOpt) {
				t.Fatalf("Length of postJson was 0 have: %q err: %v", postJson, err)
			}
			if postJson.Name != "abc" {
				http.Error(w, fmt.Sprintf("Failed send data have: %q %q", r.Body, postJson), http.StatusInternalServerError)
			}
			if !postJson.Opts.Opt1 {
				http.Error(w, fmt.Sprintf("Failed send child data have: %q %q", r.Body, postJson), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"result": "ok"}`)
		})

	tgs := httptest.NewServer(mux)

	os.Setenv(redashUrlEnv, tgs.URL)

	opt := Opt{Opt1: true, Opt2: "f"}
	opts := postOpt{Id: 123, Name: "abc", Opts: opt}
	jbuf, err := json.Marshal(opts)
	resp, err := Post("api/posttest", jbuf)
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

func TestDelete(t *testing.T) {
	beforeUrlEnv := os.Getenv(redashUrlEnv)
	defer os.Setenv(redashUrlEnv, beforeUrlEnv)

	mux := http.NewServeMux()
	mux.HandleFunc(
		"/api/deletetest",
		func(w http.ResponseWriter, r *http.Request) {
			if auth := r.Header.Get("Authorization"); !strings.Contains(auth, testApikey) {
				http.Error(w, fmt.Sprintf("Invalid Apikey %s", auth), http.StatusForbidden)
				return
			}
			if meth := r.Method; meth != http.MethodDelete {
				http.Error(w, fmt.Sprintf("Invalid Method %s", meth), http.StatusBadRequest)
				return
			}
			if id := r.URL.Query().Get("id"); id != "abc" {
				http.Error(w, fmt.Sprintf("Invalid Parameter %s", id), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"result": "ok"}`)
		})

	tgs := httptest.NewServer(mux)

	os.Setenv(redashUrlEnv, tgs.URL)

	params := map[string]string{"id": "abc"}
	resp, err := Delete("api/deletetest", params)
	if err != nil {
		t.Fatalf("Failed recieve Response from Delete. %v", err)
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
		t.Fatalf("HTTPMethod was not Delete, method: %s", resp.Body)
	}
	if want := "ok"; have.Result != want {
		t.Fatalf("resp is bad. want: %q, have: %q", want, have)
	}
}

func TestDo(t *testing.T) {
	beforeUrlEnv := os.Getenv(redashUrlEnv)
	defer os.Setenv(redashUrlEnv, beforeUrlEnv)

	mux := http.NewServeMux()
	mux.HandleFunc(
		"/api/dotest",
		func(w http.ResponseWriter, r *http.Request) {
			if auth := r.Header.Get("Authorization"); !strings.Contains(auth, testApikey) {
				http.Error(w, fmt.Sprintf("Invalid Apikey %s", auth), http.StatusForbidden)
				return
			}
			if meth := r.Method; meth != http.MethodGet {
				http.Error(w, fmt.Sprintf("Invalid Method %s", meth), http.StatusBadRequest)
				return
			}
			if id := r.URL.Query().Get("id"); id != "abc" {
				http.Error(w, fmt.Sprintf("Invalid Parameter %s", id), http.StatusInternalServerError)
			}
			w.Header().Set("Content-Type", "application/json")
			fmt.Fprintf(w, `{"result": "ok"}`)
		})

	tgs := httptest.NewServer(mux)

	os.Setenv(redashUrlEnv, tgs.URL)

	defaultOpts.Params = map[string]string{"id": "abc"}
	resp, err := Do(http.MethodGet, "api/dotest", defaultOpts)
	if err != nil {
		t.Fatalf("Failed recieve Response from Delete. %v", err)
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
		t.Fatalf("HTTPMethod was not Delete, method: %s", resp.Body)
	}
	if want := "ok"; have.Result != want {
		t.Fatalf("resp is bad. want: %q, have: %q", want, have)
	}
}

func TestRequest(t *testing.T) {
	beforeUrlEnv := os.Getenv(redashUrlEnv)
	defer os.Setenv(redashUrlEnv, beforeUrlEnv)

	os.Setenv(redashUrlEnv, testUrl)

	opt := Opt{Opt1: true, Opt2: "f"}
	opts := postOpt{Id: 123, Name: "abc", Opts: opt}
	jbuf, err := json.Marshal(opts)
	defaultOpts.Body = bytes.NewReader(jbuf)
	req, err := Request(http.MethodPost, "api/reqtest", defaultOpts)

	if auth := req.Header.Get("Authorization"); !strings.Contains(auth, testApikey) {
		t.Fatalf("Invalid Apikey %s", auth)
	}
	if meth := req.Method; meth != http.MethodPost {
		t.Fatalf("Invalid Method %s", meth)
	}
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		t.Fatalf("Failed ReadAll have: %q err: %v", req.Body, err)
	}
	var postJson postOpt
	if err := json.Unmarshal(buf, &postJson); err != nil {
		t.Fatalf("Failed Json Decode have: %q err: %v", postJson, err)
	}
	if postJson.Name != "abc" {
		t.Fatalf("Failed prepare data have: %q %q", req.Body, postJson)
	}
}

type originalClientData struct {}

var (
	origTestKey = "originalApikey"
	origTestUrl = "http://original.com"
)

func (originalClientData) Apikey() (apikey string, err error) {
	return origTestKey, nil
}

func (originalClientData) Url() (u *url.URL, err error) {
	return url.Parse(origTestUrl)
}

func (originalClientData) HTTPClient() *http.Client {
	return &http.Client{}
}

func (originalClientData) DefaultOpts() *Options {
	return &defaultOpts
}


func TestOriginalClient(t *testing.T) {

	client := originalClientData{}

	u, err := client.Url()
	if err != nil {
		t.Error(err)
	}
	if u.String() != origTestUrl {
		t.Fatalf("Url is bad. want: %q, have: %q", origTestKey, u.String())
	}

	apikey, err := client.Apikey()
	if err != nil {
		t.Error(err)
	}
	if apikey != origTestKey {
		t.Fatalf("Failed get Apikey want %q have %q", origTestKey, apikey)
	}

	hc := client.HTTPClient()
	nc := &http.Client{}
	if !reflect.DeepEqual(hc, nc) {
		t.Fatalf("HTTPClient is bad. want: %q, have: %q", nc, hc)
	}

	if !reflect.DeepEqual(client.DefaultOpts(), &defaultOpts) {
		t.Fatalf("DefaultOptions is bad. want: %q, have: %q", &defaultOpts, client.DefaultOpts())
	}

}

func TestGetInter(t *testing.T) {}

func TestPostInter(t *testing.T) {}

func TestDeleteInter(t *testing.T) {}

func TestDoInter(t *testing.T) {}

func TestRequestInter(t *testing.T) {}
