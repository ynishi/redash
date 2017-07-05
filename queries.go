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

var Queries = &QueriesS{DefaultClient}

type Querieser interface {
	Queries(string) string
}

type QueriesS struct {
	Client Interface
}

func (q QueriesS) Queries(s string) (rs string) {
	return "/api/queries/" + s
}

type FormatQuery struct {
	Query string `json:"query"`
}

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

type NewQuery struct {
	DataSourceId int               `json:"data_source_id"`
	Query        string            `json:"query"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Schedule     string            `json:"schedule"`
	Options      map[string]string `json:"options"`
}

func (qs QueriesS) PostFormat(sql string) (r io.Reader, err error) {
	resp, err := PostInter(qs.Client, qs.Queries("format"), []byte(fmt.Sprintf(`{"query":"%s"}`, sql)))
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

func (qs QueriesS) GetSearch(q string) (r io.Reader, err error) {
	params := map[string]string{"q": q}
	resp, err := GetInter(qs.Client, qs.Queries("search"), params)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

func (qs QueriesS) GetRecent() (r io.Reader, err error) {
	resp, err := GetInter(qs.Client, qs.Queries("recent"), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

func (qs QueriesS) GetMy(pageSize, page int) (r io.Reader, err error) {
	params := map[string]string{"page_size": strconv.Itoa(pageSize), "page": strconv.Itoa(page)}
	resp, err := GetInter(qs.Client, qs.Queries("my"), params)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

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

func (qs QueriesS) GetQuery(pageSize, page int) (r io.Reader, err error) {
	params := map[string]string{"page_size": strconv.Itoa(pageSize), "page": strconv.Itoa(page)}
	resp, err := GetInter(qs.Client, qs.Queries(""), params)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

func (qs QueriesS) PostRefresh(queryId int) (r io.Reader, err error) {
	resp, err := PostInter(qs.Client, qs.Queries(fmt.Sprintf("%s/refresh", strconv.Itoa(queryId))), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

func (qs QueriesS) PostFork(queryId int) (r io.Reader, err error) {
	resp, err := PostInter(qs.Client, qs.Queries(fmt.Sprintf("%s/fork", strconv.Itoa(queryId))), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

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

func (qs QueriesS) DeleteQuery(queryId int) (r io.Reader, err error) {
	resp, err := DeleteInter(qs.Client, qs.Queries(strconv.Itoa(queryId)), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

func (qs QueriesS) GetQueryId(queryId int) (r io.Reader, err error) {
	resp, err := GetInter(qs.Client, qs.Queries(strconv.Itoa(queryId)), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

func (qs QueriesS) PostQueryResult(query string, queryId, maxAge, dataSourceId int) (r io.Reader, err error) {
	resp, err := PostInter(qs.Client, "api/query_results", []byte(fmt.Sprintf(`{"query":"%s","query_id":%d,"max_age":"%s","data_sourece_id":%d}`)))
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

//GET /api/queries/(query_id)/results/(query_result_id).(filetype)
func (qs QueriesS) GetResultsId(queryId, queryResultId int, filetype string) (r io.Reader, err error) {
	resp, err := GetInter(qs.Client, qs.Queries(fmt.Sprintf("%s/results/%s.%s", strconv.Itoa(queryId), strconv.Itoa(queryResultId), filetype)), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

//GET /api/queries/(query_id)/results.(filetype)
func (qs QueriesS) PostResults(queryId int, filetype string) (r io.Reader, err error) {
	resp, err := PostInter(qs.Client, qs.Queries(fmt.Sprintf("%s/results.%s", strconv.Itoa(queryId), filetype)), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

//GET /api/query_results/(query_result_id)
func (qs QueriesS) GetQueryResults(queryId, queryResultId int, filetype string) (r io.Reader, err error) {
	params := map[string]string{"query_id": strconv.Itoa(queryId), "query_result_id": strconv.Itoa(queryResultId), "filetype": filetype}
	resp, err := GetInter(qs.Client, fmt.Sprintf("query_results/%s", queryResultId), params)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

//DELETE /api/jobs/(job_id)
func (qs QueriesS) DeleteJog(jobId int) (r io.Reader, err error) {
	resp, err := DeleteInter(qs.Client, fmt.Sprintf("api/jobs/%s", strconv.Itoa(jobId)), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}

//GET /api/jobs/(job_id)
func (qs QueriesS) GetJob(jobId int) (r io.Reader, err error) {
	resp, err := GetInter(qs.Client, fmt.Sprintf("api/jobs/%s", strconv.Itoa(jobId)), nil)
	if err != nil {
		return nil, err
	} else {
		return resp.Body, nil
	}
}
