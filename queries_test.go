package redash

import (
	"net/http"
	"testing"

	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"strings"
)

type muxVal struct {
	path       string
	getResp    string
	postResp   string
	deleteResp string
}

var (
	tgs     *httptest.Server
	muxData = []muxVal{
		{"queries/format", "", `{"query": "select 1;"}`, ""},
		{"queries/search", `[{"name": "queryTitle1"}]`, "", ""},
		{"queries/recent", `[{"name": "queryTitle1"}]`, "", ""},
		{"queries/my", `[{"name": "queryTitle1"}]`, "", ""},
		{"queries", `{"id":1}`, `{"id": 1}`, ""},
		{"queries/1", `{"id": 1}`, `{"id": 1}`, `{"id": 1}`},
		{"queries/1/refresh", "", `{"id": 1}`, ""},
		{"queries/1/fork", "", `{"id": 1}`, ""},
		{"queries/1/results/1.json", `{"id": 1}`, "", ""},
		{"queryes/1/results.json", `{"id": 1}`, "", ""},
		{"query_results", "", `{"id": 1}`, ""},
		{"query_results/1", `{"id": 1}`, "", ""},
		{"jobs/1", `{"id": 1}`, "", `{"id": 1}`},
	}
)

func init() {

	mux := http.NewServeMux()
	for _, d := range muxData {
		ep := fmt.Sprintf("/api/%s", d.path)
		rg := d.getResp
		rp := d.postResp
		rd := d.deleteResp
		var rs string
		mux.HandleFunc(
			ep,
			func(w http.ResponseWriter, r *http.Request) {
				if auth := r.Header.Get("Authorization"); !strings.Contains(auth, mockApikey) {
					http.Error(w, fmt.Sprintf("Invalid Apikey %s", auth), http.StatusForbidden)
					return
				}
				switch r.Method {
				case http.MethodGet:
					if rg == "" {
						http.Error(w, "Method Get not supported.", http.StatusBadRequest)
						return
					} else {
						rs = rg
					}
				case http.MethodPost:
					if rp == "" {
						http.Error(w, "Method Post not supported.", http.StatusBadRequest)
						return
					} else {
						rs = rp
					}
				case http.MethodDelete:
					if rd == "" {
						http.Error(w, "Method Delete not supported.", http.StatusBadRequest)
						return
					} else {
						rs = rd
					}
				}
				w.Header().Set("Content-Type", "application/json")
				fmt.Fprintf(w, rs)
			})
	}
	tgs = httptest.NewServer(mux)

	mockClient.MockUrl = tgs.URL
	Queries.Client = mockClient
}

func TestQueriesPostFormat(t *testing.T) {

	sql := "select 1;"

	r, err := Queries.PostFormat(sql)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	var formated FormatQuery
	err = json.Unmarshal(buf, &formated)
	if err != nil {
		t.Error(err)
	}
	if formated.Query != sql {
		t.Fatalf("Query not matched, want: %q, have: %q", sql, formated.Query)
	}
}

func TestQueriesGetSearch(t *testing.T) {

	q := "queryTitle1"
	r, err := Queries.GetSearch(q)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	var responsed []ResponseQuery

	err = json.Unmarshal(buf, &responsed)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(responsed[0].Name, q) {
		t.Fatalf("Query Title not matched, want: %q, have: %q", q, responsed[0].Name)
	}
}
