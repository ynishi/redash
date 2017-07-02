// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.
package redash

import (
	"fmt"
	"io"
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
