// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package redash

import (
	"encoding/json"
	"fmt"
	"io"
	"strconv"
)

// Default Queries
var Queries = &QueriesS{DefaultClient}

// interface of make Queries endpoint.
type Querieser interface {
	Queries(string) string
}

// Default struct for queries.
type QueriesS struct {
	Client Interface
}

// Default implement of Queries.
func (q QueriesS) Queries(s string) (rs string) {
	return "/api/queries/" + s
}

// Wrap Redash format query.
type FormatQuery struct {
	Query string `json:"query"`
}

// Wrap Redash paging response query.
type PagingResponseQuery struct {
	Count    int             `json:"count"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
	Results  []ResponseQuery `json:"results"`
}

// Wrap Redash response query.
type ResponseQuery struct {
	Id                int     `json:"id"`
	LatestQueryDataId int     `json:"latest_query_data_id"`
	Name              string  `json:"name"`
	Description       string  `json:"description"`
	Query             string  `json:"query"`
	QueryHash         string  `json:"query_hash"`
	Schedule          string  `json:"schedule"`
	ApiKey            string  `json:"api_key"`
	IsArchived        bool    `json:"is_archived"`
	IsDraft           bool    `json:"is_draft"`
	UpdatedAt         string  `json:"updated_at"`
	CreatedAt         string  `json:"created_at"`
	DataSourceId      int     `json:"data_source_id"`
	Options           Options `json:"options"`
	Version           int     `json:"version"`
	UserId            int     `json:"user_id"`
	LastModifiedById  int     `json:"last_modified_by_id"`
	RetrivedAt        string  `json:"retrieved_at"`
	Runtime           int     `json:"runtime"`
}

// Wrap Redash new query.
type NewQuery struct {
	DataSourceId int               `json:"data_source_id"`
	Query        string            `json:"query"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Schedule     string            `json:"schedule"`
	Options      map[string]string `json:"options"`
}

// Wrap Redash row for result data.
type Row struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// Wrap Redash column for result data.
type Column struct {
	FriendlyName string `json:"friendly_name"`
	Type         string `json:"type"`
	Name         string `json:"name"`
}

// Wrap Redash result data.
type ResultData struct {
	Rows    []Row    `json:"rows"`
	Columns []Column `json:"columns"`
}

// Wrap Redash query result.
type QueryResult struct {
	RetrievedAt  string     `json:"retrieved_at"`
	QueryHash    string     `json:"query_hash"`
	Query        string     `json:"query"`
	Runtime      float64    `json:"runtime"`
	Data         ResultData `json:"data"`
	Id           int        `json:"id"`
	DataSourceId int        `json:"data_source_id"`
}

// Wrap Redash result.
type Result struct {
	QueryResult QueryResult `json:"query_result"`
}

// Wrap Redash job detail.
type JobInner struct {
	Status        int    `json:"status"`
	Error         string `json:"error"`
	Id            string `json:"id"`
	QueryResultId int    `json:"query_result_id"`
	Updated_at    int    `json:"updated_at"`
}

// Wrap Redash job result.
type Job struct {
	Job JobInner `josn:"job"`
}

// Wrap Redash api POST format.
func (qs QueriesS) PostFormat(sql string) (r io.Reader, err error) {
	resp, err := PostInter(qs.Client, qs.Queries("format"), []byte(fmt.Sprintf(`{"query":"%s"}`, sql)))
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api GET search.
func (qs QueriesS) GetSearch(q string) (r io.Reader, err error) {
	params := map[string]string{"q": q}
	resp, err := GetInter(qs.Client, qs.Queries("search"), params)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api GET recent.
func (qs QueriesS) GetRecent() (r io.Reader, err error) {
	resp, err := GetInter(qs.Client, qs.Queries("recent"), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api GET my.
func (qs QueriesS) GetMy(pageSize, page int) (r io.Reader, err error) {
	params := map[string]string{"page_size": strconv.Itoa(pageSize), "page": strconv.Itoa(page)}
	resp, err := GetInter(qs.Client, qs.Queries("my"), params)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api POST queries.
func (qs QueriesS) PostQuery(newQuery NewQuery) (res io.Reader, err error) {
	newQueryBuf, err := json.Marshal(newQuery)
	if err != nil {
		return nil, err
	}
	resp, err := PostInter(qs.Client, qs.Queries(""), newQueryBuf)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api GET queries.
func (qs QueriesS) GetQuery(pageSize, page int) (r io.Reader, err error) {
	params := map[string]string{"page_size": strconv.Itoa(pageSize), "page": strconv.Itoa(page)}
	resp, err := GetInter(qs.Client, qs.Queries(""), params)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api POST refresh.
func (qs QueriesS) PostRefresh(queryId int) (r io.Reader, err error) {
	resp, err := PostInter(qs.Client, qs.Queries(fmt.Sprintf("%d/refresh", queryId)), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api POST fork.
func (qs QueriesS) PostFork(queryId int) (r io.Reader, err error) {
	resp, err := PostInter(qs.Client, qs.Queries(fmt.Sprintf("%d/fork", queryId)), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api POST queries.
func (qs QueriesS) PostQueryId(queryId int, newQuery NewQuery) (r io.Reader, err error) {
	newQueryBuf, err := json.Marshal(newQuery)
	if err != nil {
		return nil, err
	}
	resp, err := PostInter(qs.Client, qs.Queries(strconv.Itoa(queryId)), newQueryBuf)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api DELETE queries.
func (qs QueriesS) DeleteQuery(queryId int) (r io.Reader, err error) {
	resp, err := DeleteInter(qs.Client, qs.Queries(strconv.Itoa(queryId)), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api GET queries/${query id}.
func (qs QueriesS) GetQueryId(queryId int) (r io.Reader, err error) {
	resp, err := GetInter(qs.Client, qs.Queries(strconv.Itoa(queryId)), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api POST query_results.
func (qs QueriesS) PostQueryResult(query string, maxAge, dataSourceId int) (r io.Reader, err error) {
	resp, err := PostInter(qs.Client, "/api/query_results", []byte(fmt.Sprintf(`{"query":"%s","max_age":%d,"data_sourece_id":%d}`, query, maxAge, dataSourceId)))
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api GET ${query id}/results/${query resut id}.${filetype}
func (qs QueriesS) GetResultsById(queryId, queryResultId int, filetype string) (r io.Reader, err error) {
	resp, err := GetInter(qs.Client, qs.Queries(fmt.Sprintf("%d/results/%d.%s", queryId, queryResultId, filetype)), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api GET ${query id}/results.${filetype}.
func (qs QueriesS) GetResultsByQueryId(queryId int, filetype string) (r io.Reader, err error) {
	resp, err := GetInter(qs.Client, qs.Queries(fmt.Sprintf("%d/results.%s", queryId, filetype)), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api GET querie_results.
func (qs QueriesS) GetQueryResults(queryResultId int) (r io.Reader, err error) {
	resp, err := GetInter(qs.Client, fmt.Sprintf("/api/query_results/%d", queryResultId), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api DELETE jobs.
func (qs QueriesS) DeleteJog(jobId string) (r io.Reader, err error) {
	resp, err := DeleteInter(qs.Client, fmt.Sprintf("/api/jobs/%s", jobId), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

// Wrap Redash api GET jobs.
func (qs QueriesS) GetJob(jobId string) (r io.Reader, err error) {
	resp, err := GetInter(qs.Client, fmt.Sprintf("/api/jobs/%s", jobId), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}
