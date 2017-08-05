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

const fromatResp = `{
  "query": "SELECT *\nFROM dual;"
}`

const queryResp = `{
  "last_modified_by_id": 1,
  "latest_query_data_id": 2,
  "schedule": null,
  "is_archived": false,
  "updated_at": "2017-07-16T10:52:26.541613+00:00",
  "user": {
    "auth_type": "password",
    "created_at": "2017-07-16T10:15:31.897134+00:00",
    "name": "user1",
    "gravatar_url": "",
    "updated_at": "2017-07-16T10:15:31.897134+00:00",
    "id": 1,
    "groups": [
      1
    ],
    "email": "user1@example.com"
  },
  "query": "select * from hello;",
  "is_draft": false,
  "id": 1,
  "description": null,
  "name": "helloQuery",
  "created_at": "2017-07-16T10:43:33.399535+00:00",
  "version": 1,
  "query_hash": "31d0721311ab1bcccc07504e61aa20dc",
  "api_key": "abcdef",
  "options": {
    "parameters": []
  },
  "data_source_id": 1
}`

var queryResps = fmt.Sprintf(`[%s]`, queryResp)

var pagingResp = fmt.Sprintf(`{
  "count": 1,
  "page": 1,
  "page_size": 20,
  "results": %s
}`, queryResps)

const queryResultResp = `{
  "query_result": {
    "retrieved_at": "2017-07-16T11:49:35.033971+00:00",
    "query_hash": "31d0721311ab1bcccc07504e61aa20dc",
    "query": "select * from hello;",
    "runtime": 0.00354194641113281,
    "data": {
      "rows": [
        {
          "id": 1,
          "name": "test1"
        },
        {
          "id": 2,
          "name": "test2"
        }
      ],
      "columns": [
        {
          "friendly_name": "id",
          "type": "integer",
          "name": "id"
        },
        {
          "friendly_name": "name",
          "type": null,
          "name": "name"
        }
      ]
    },
    "id": 2,
    "data_source_id": 1
  }
}`

const newQueryResp = `{
  "latest_query_data_id": null,
  "schedule": null,
  "is_archived": false,
  "updated_at": "2017-07-16T13:29:45.356947+00:00",
  "user": {
    "auth_type": "password",
    "created_at": "2017-07-16T10:15:31.897134+00:00",
    "name": "admin",
    "gravatar_url": "",
    "updated_at": "2017-07-16T10:15:31.897134+00:00",
    "id": 1,
    "groups": [
      1,
      2
    ],
    "email": "user1@example.com"
  },
  "query": "select * from hello;",
  "is_draft": true,
  "id": 1,
  "description": null,
  "name": "api",
  "created_at": "2017-07-16T13:29:45.356947+00:00",
  "last_modified_by": {
    "auth_type": "password",
    "created_at": "2017-07-16T10:15:31.897134+00:00",
    "name": "admin",
    "gravatar_url": "",
    "updated_at": "2017-07-16T10:15:31.897134+00:00",
    "id": 1,
    "groups": [
      1,
      2
    ],
    "email": "user1@example.com"
  },
  "version": 1,
  "query_hash": "914a74181b749b366dfaebf7aaf52164",
  "api_key": "abcdef",
  "options": {},
  "data_source_id": 1
}`

const updateQueryResp = `{
  "latest_query_data_id": null,
  "schedule": null,
  "is_archived": false,
  "updated_at": "2017-07-16T13:29:45.356947+00:00",
  "user": {
    "auth_type": "password",
    "created_at": "2017-07-16T10:15:31.897134+00:00",
    "name": "admin",
    "gravatar_url": "",
    "updated_at": "2017-07-16T10:15:31.897134+00:00",
    "id": 1,
    "groups": [
      1,
      2
    ],
    "email": "user1@example.com"
  },
  "query": "select * from hello;",
  "is_draft": true,
  "id": 1,
  "description": null,
  "name": "api2",
  "created_at": "2017-07-16T13:29:45.356947+00:00",
  "last_modified_by": {
    "auth_type": "password",
    "created_at": "2017-07-16T10:15:31.897134+00:00",
    "name": "admin",
    "gravatar_url": "",
    "updated_at": "2017-07-16T10:15:31.897134+00:00",
    "id": 1,
    "groups": [
      1,
      2
    ],
    "email": "user1@example.com"
  },
  "version": 1,
  "query_hash": "914a74181b749b366dfaebf7aaf52164",
  "api_key": "abcdef",
  "options": {},
  "data_source_id": 1
}`

const forkQueryResp = `{
  "latest_query_data_id": null,
  "schedule": null,
  "is_archived": false,
  "updated_at": "2017-07-16T13:29:45.356947+00:00",
  "user": {
    "auth_type": "password",
    "created_at": "2017-07-16T10:15:31.897134+00:00",
    "name": "admin",
    "gravatar_url": "",
    "updated_at": "2017-07-16T10:15:31.897134+00:00",
    "id": 1,
    "groups": [
      1,
      2
    ],
    "email": "user1@example.com"
  },
  "query": "select * from hello;",
  "is_draft": true,
  "id": 2,
  "description": null,
  "name": "Copy of (#1) api",
  "created_at": "2017-07-16T13:29:45.356947+00:00",
  "last_modified_by": {
    "auth_type": "password",
    "created_at": "2017-07-16T10:15:31.897134+00:00",
    "name": "admin",
    "gravatar_url": "",
    "updated_at": "2017-07-16T10:15:31.897134+00:00",
    "id": 1,
    "groups": [
      1,
      2
    ],
    "email": "user1@example.com"
  },
  "version": 1,
  "query_hash": "914a74181b749b366dfaebf7aaf52164",
  "api_key": "abcdef",
  "options": {},
  "data_source_id": 1
}`

const jobResp = `{
  "job": {
    "status": 2,
    "error": "",
    "id": "d856637d-9387-4874-a944-9c93ac45c688",
    "query_result_id": null,
    "updated_at": 0
  }
}`

const jobRet = `{
  "job": {
    "status": 3,
    "error": "",
    "id": "d856637d-9387-4874-a944-9c93ac45c688",
    "query_result_id": 1,
    "updated_at": 0
  }
}`

var (
	tgs     *httptest.Server
	muxData = []muxVal{
		{"queries/format", "", fromatResp, ""},
		{"queries/search", queryResps, "", ""},
		{"queries/recent", queryResps, "", ""},
		{"queries/my", pagingResp, "", ""},
		{"queries", pagingResp, newQueryResp, ""},
		{"queries/1", queryResp, updateQueryResp, "null"},
		{"queries/1/refresh", "", jobResp, ""},
		{"queries/1/fork", "", forkQueryResp, ""},
		{"queries/1/results/2.json", queryResultResp, "", ""},
		{"queries/1/results.json", queryResultResp, "", ""},
		{"query_results", "", jobResp, ""},
		{"query_results/2", queryResultResp, "", ""},
		{"jobs/d856637d-9387-4874-a944-9c93ac45c688", jobRet, "", "null"},
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

func TestPostFormat(t *testing.T) {

	sql := "select 1 from dual;"

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
	var wanted FormatQuery
	_ = json.Unmarshal([]byte(fromatResp), &wanted)
	if formated.Query != wanted.Query {
		t.Fatalf("Query not matched,\n want: %q,\n have: %q\n", wanted.Query, formated.Query)
	}
}

func TestGetSearch(t *testing.T) {

	q := "hello"

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
		t.Fatalf("Query Title not matched,\n want: %q,\n have: %q\n", q, responsed[0].Name)
	}
}

func TestGetRecent(t *testing.T) {

	name := "helloQuery"

	r, err := Queries.GetRecent()
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	var responsed []ResponseQuery
	err = json.Unmarshal(buf, &responsed)
	if err != nil {
		t.Error(err)
	}
	if !strings.Contains(name, responsed[0].Name) {
		t.Fatalf("Query Title not matched,\n want: %q,\n have: %q\n", name, responsed[0].Name)
	}
}

func TestGetMy(t *testing.T) {

	pageSize := 20
	page := 1

	r, err := Queries.GetMy(pageSize, page)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	var responsed PagingResponseQuery
	err = json.Unmarshal(buf, &responsed)
	if err != nil {
		t.Error(err)
	}
	if responsed.PageSize != pageSize {
		t.Fatalf("PageSize was not match,\n want: %q,\n have: %q\n", pageSize, responsed.PageSize)
	}
	if responsed.Page != page {
		t.Fatalf("Page was not match,\n want: %q,\n have: %q\n", page, responsed.Page)
	}
}

func TestPostQuery(t *testing.T) {

	dataSourceId := 1
	query := "select * from hello;"
	name := "api"
	newQuery := &NewQuery{
		DataSourceId: dataSourceId,
		Query:        query,
		Name:         name,
	}

	r, err := Queries.PostQuery(*newQuery)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	var responsed ResponseQuery
	err = json.Unmarshal(buf, &responsed)
	if err != nil {
		t.Error(err)
	}
	if responsed.DataSourceId != dataSourceId {
		t.Fatalf("DataSourceId is not match,\n want: %q,\n have: %q\n", dataSourceId, responsed.DataSourceId)
	}
	if responsed.Query != query {
		t.Fatalf("Query is not match,\n want: %q,\n have: %q\n", query, responsed.Query)
	}
	if responsed.Name != name {
		t.Fatalf("Name is not match,\n want: %q,\n have: %q\n", name, responsed.Name)
	}
}

func TestGetQuery(t *testing.T) {

	pageSize := 20
	page := 1
	r, err := Queries.GetQuery(pageSize, page)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	var responsed PagingResponseQuery

	err = json.Unmarshal(buf, &responsed)
	if err != nil {
		t.Error(err)
	}
	if responsed.PageSize != pageSize {
		t.Fatalf("PageSize was not match,\n want: %q,\n have: %q\n", pageSize, responsed.PageSize)
	}
	if responsed.Page != page {
		t.Fatalf("Page was not match,\n want: %q,\n have: %q\n", page, responsed.Page)
	}
}

func TestPostRefresh(t *testing.T) {

	queryId := 1
	r, err := Queries.PostRefresh(queryId)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	var job Job

	err = json.Unmarshal(buf, &job)
	if err != nil {
		t.Error(err)
	}
	if job.Job.Status != 2 {
		t.Fatalf("Job status is not match,\n want: %q,\n have: %q\n", 2, job.Job.Status)
	}
	if job.Job.QueryResultId != 0 {
		t.Fatalf("Query result id is not match,\n want: %q,\n have: %q\n", "", job.Job.QueryResultId)
	}
}

func TestPostFork(t *testing.T) {

	queryId := 1
	forkedName := "Copy of (#1) api"

	r, err := Queries.PostFork(queryId)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	var responsed ResponseQuery

	err = json.Unmarshal(buf, &responsed)
	if err != nil {
		t.Error(err)
	}
	if responsed.Id == queryId {
		t.Fatalf("Query id is not updated,\n want: != %q,\n have: %q\n", queryId, responsed.Id)
	}
	if responsed.Name != forkedName {
		t.Fatalf("Name is not match,\n want: %q,\n have: %q\n", forkedName, responsed.Name)
	}
}

func TestPostQueryId(t *testing.T) {

	queryId := 1
	dataSourceId := 1
	query := "select * from hello;"
	name := "api2"
	updateQuery := &NewQuery{
		DataSourceId: dataSourceId,
		Query:        query,
		Name:         name,
	}

	r, err := Queries.PostQueryId(queryId, *updateQuery)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	var responsed ResponseQuery
	err = json.Unmarshal(buf, &responsed)
	if err != nil {
		t.Error(err)
	}
	if responsed.Id != queryId {
		t.Fatalf("Query id is not match,\n want: %q,\n have: %q\n", queryId, responsed.Id)
	}
	if responsed.Query != query {
		t.Fatalf("Query is not match,\n want: %q,\n have: %q\n", query, responsed.Query)
	}
	if responsed.Name != name {
		t.Fatalf("Name is not match,\n want: %q,\n have: %q\n", name, responsed.Name)
	}
}

func TestDeleteQuery(t *testing.T) {

	queryId := 1

	r, err := Queries.DeleteQuery(queryId)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error(err)
	}
	if string(buf) != "null" {
		t.Fatalf("Resq is invalid,\n want: %q,\n have: %q\n", "null", string(buf))
	}
}

func TestGetQueryId(t *testing.T) {

	queryId := 1
	name := "helloQuery"

	r, err := Queries.GetQueryId(queryId)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	var responsed ResponseQuery
	err = json.Unmarshal(buf, &responsed)
	if err != nil {
		t.Fatalf("Format is not json,\n err: %q,\n have: %q\n", err, buf)
	}
	if responsed.Id != queryId {
		t.Fatalf("Query id is not match,\n want: %q,\n have: %q\n", queryId, responsed.Id)
	}
	if responsed.Name != name {
		t.Fatalf("Query name is not match,\n want: %q,\n have: %q\n", name, responsed.Name)
	}
}

func TestPostQueryResult(t *testing.T) {

	dataSourceId := 1
	maxAge := 0
	query := "select * from hello;"

	r, err := Queries.PostQueryResult(query, maxAge, dataSourceId)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	var respJob Job

	err = json.Unmarshal(buf, &respJob)
	if err != nil {
		t.Error(err)
	}
	if respJob.Job.Status != 2 {
		t.Fatalf("Job status is error,\n want: %q,\n have: %q\n", 2, respJob.Job.Status)
	}
	if respJob.Job.QueryResultId != 0 {
		t.Fatalf("QueryResultId is not null,\n want: %q,\n have: %q\n", "null", respJob.Job.QueryResultId)
	}
	if respJob.Job.Error != "" {
		t.Fatalf("Error is not empty,\n want: %q,\n have: %q\n", "", respJob.Job.Error)
	}
}

func TestGetResultsById(t *testing.T) {

	queryId := 1
	queryResultId := 2
	filetype := "json"

	r, err := Queries.GetResultsById(queryId, queryResultId, filetype)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	var results Result
	err = json.Unmarshal(buf, &results)
	if err != nil {
		t.Fatalf("Format is not json,\n err: %q,\n have: %q\n", err, buf)
	}
	if results.QueryResult.Id != queryResultId {
		t.Fatalf("Query result id is not match,\n want: %q,\n have: %q\n", queryResultId, results.QueryResult.Id)
	}
	if len(results.QueryResult.Data.Rows) == 0 {
		t.Fatalf("Rows num is bad,\n want: %q,\n have: %q\n", 2, len(results.QueryResult.Data.Rows))
	}
}

func TestGetResultsByQueryId(t *testing.T) {

	queryId := 1
	filetype := "json"

	r, err := Queries.GetResultsByQueryId(queryId, filetype)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	var results Result
	err = json.Unmarshal(buf, &results)
	if err != nil {
		t.Fatalf("Format is not json,\n err: %q,\n hav	e: %q\n", err, buf)
	}
	if len(results.QueryResult.Data.Rows) == 0 {
		t.Fatalf("Rows num is bad,\n want: %q,\n have: %q\n", 2, len(results.QueryResult.Data.Rows))
	}
}

func TestGetQueryResults(t *testing.T) {

	queryResultId := 2

	r, err := Queries.GetQueryResults(queryResultId)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	var results Result
	err = json.Unmarshal(buf, &results)
	if err != nil {
		t.Fatalf("Format is not json,\n err: %q,\n have: %q\n", err, buf)
	}
	if results.QueryResult.Id != queryResultId {
		t.Fatalf("Query result id is not match,\n want: %q,\n have: %q\n", queryResultId, results.QueryResult.Id)
	}
	if len(results.QueryResult.Data.Rows) == 0 {
		t.Fatalf("Rows num is bad,\n want: %q,\n have: %q\n", 2, len(results.QueryResult.Data.Rows))
	}
}

func TestDeleteJog(t *testing.T) {

	jobId := "d856637d-9387-4874-a944-9c93ac45c688"

	r, err := Queries.DeleteJog(jobId)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		t.Error(err)
	}
	if string(buf) != "null" {
		t.Fatalf("Resq is invalid,\n want: %q,\n have: %q\n", "null", string(buf))
	}
}

func TestGetJob(t *testing.T) {

	jobId := "d856637d-9387-4874-a944-9c93ac45c688"

	r, err := Queries.GetJob(jobId)
	if err != nil {
		t.Error(err)
	}
	buf, err := ioutil.ReadAll(r)
	var job Job
	err = json.Unmarshal(buf, &job)
	if err != nil {
		t.Fatalf("Format is not json,\n err: %q,\n have: %q\n", err, buf)
	}
	if job.Job.Id != jobId {
		t.Fatalf("Job id is not match,\n want: %q,\n have: %q\n", jobId, job.Job.Id)
	}
	if job.Job.Error != "" {
		t.Fatalf("Error is not empty,\n want: %q,\n have: %q\n", "", job.Job.Error)
	}
}
