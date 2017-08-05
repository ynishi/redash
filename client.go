// Copyright 2017 Yutaka Nishimura. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

/*
   The redash package implements a simple client and wrapper library for
   Redash REST api.

   GET/POST/DELETE is accepted.

   Original client is made by implement Interface.

   Summary of use case.

   Case 1:
     Get/Post/Delete directory.
   Case 2:
     implement your Interface and
     GetInter/PostInter/DeleteInter with new Client.
   Case 3:
     Queries.GetQuery/(other func in Queries)

*/
package redash

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
)

const (
	redashUrlEnv    = "REDASH_URL"
	redashApikeyEnv = "REDASH_APIKEY"
)

var (
	repository = "https://github.com/ynishi/redash"
	ua         = fmt.Sprintf("RedashGoClient/0.1 (+%s; %s)",
		repository, runtime.Version())
	defaultPostHeader = map[string]string{
		"Content-Type": "application/json",
	}
	DefaultClient = NewDefaultClient()
)

func defaultOpts() *Options {
	return &Options{
		Params: make(map[string]string),
		Header: map[string]string{
			"User-Agent": ua,
		},
		Body: nil,
	}
}

// Options is option value container.
type Options struct {
	Params map[string]string
	Header map[string]string
	Body   io.Reader
}

// Url is Redash server's endpoint.
type Urler interface {
	Url() (*url.URL, error)
}

// Apikey is Redash Apikey to connect primary.
type Apikeyer interface {
	Apikey() (string, error)
}

// HTTPClient is HTTP client to do request.
type HTTPClienter interface {
	HTTPClient() *http.Client
}

// DefaultOpts is default options for request.
type DefaultOptser interface {
	DefaultOpts() *Options
}

// A type, for original client, that sufisfies redash.Interface can
// use GetInter, PostInter, DeleteInter methods.
// If want to change apikey management, Httpclient and so on,
// just implement methods of this Interface.
type Interface interface {
	Urler
	Apikeyer
	HTTPClienter
	DefaultOptser
}

// GetInter do Redash GET with Interface and return result.
func GetInter(client Interface, sub string, params map[string]string) (resp *http.Response, err error) {
	opts := client.DefaultOpts()
	for key, value := range params {
		opts.Params[key] = value
	}
	return DoInter(client, http.MethodGet, sub, opts)
}

// PostInter do Redash POST with Interface and return result.
func PostInter(client Interface, sub string, jsonBody []byte) (resp *http.Response, err error) {
	opts := client.DefaultOpts()
	for key, value := range defaultPostHeader {
		opts.Header[key] = value
	}
	opts.Body = bytes.NewReader(jsonBody)
	return DoInter(client, http.MethodPost, sub, opts)
}

// DeleteInter do Redash DELETE with Interface and return result.
func DeleteInter(client Interface, sub string, params map[string]string) (resp *http.Response, err error) {
	opts := client.DefaultOpts()
	for key, value := range params {
		opts.Params[key] = value
	}
	return DoInter(client, http.MethodDelete, sub, opts)
}

// DoInter do Redash apis with Interface and return result.
func DoInter(client Interface, method, sub string, opts *Options) (resp *http.Response, err error) {
	log.Printf("[INFO] do: %s %s", method, sub)
	req, err := RequestInter(client, method, sub, opts)
	if err != nil {
		return nil, err
	}
	return client.HTTPClient().Do(req)
}

// RequestInter make request with Interface.
func RequestInter(client Interface, method, sub string, opts *Options) (req *http.Request, err error) {
	u, err := client.Url()
	if err != nil {
		return nil, err
	}
	u.Path = path.Join(u.Path, sub)
	values := url.Values{}
	for key, value := range opts.Params {
		values.Add(key, value)
	}
	req, err = http.NewRequest(method, u.String(), opts.Body)
	req.URL.RawQuery = values.Encode()
	if err != nil {
		return nil, err
	}
	apikey, err := client.Apikey()
	if err != nil {
		return nil, err
	}
	opts.Header["Authorization"] = "Key " + apikey
	for key, value := range opts.Header {
		req.Header.Set(key, value)
	}
	return req, nil
}

// Get do Redash api GET and return result.
func Get(sub string, params map[string]string) (resp *http.Response, err error) {
	return GetInter(DefaultClient, sub, params)
}

// Post do Redash api POST and return result.
func Post(sub string, jsonBody []byte) (resp *http.Response, err error) {
	return PostInter(DefaultClient, sub, jsonBody)
}

// Delete do Redash api DELETE and return result.
func Delete(sub string, params map[string]string) (resp *http.Response, err error) {
	return DeleteInter(DefaultClient, sub, params)
}

// Do do Redash apis and return result.
func Do(method, sub string, opts *Options) (resp *http.Response, err error) {
	return DoInter(DefaultClient, method, sub, opts)
}

// Request make http.Request for Redash.
func Request(method, sub string, opts *Options) (req *http.Request, err error) {
	return RequestInter(DefaultClient, method, sub, opts)
}

// Default implemet of client include Logger.
type ClientData struct {
	*log.Logger
}

// Default implement of client. This is provided as DefaultClient.
type DefaultClientData struct {
	ClientData
	apikey string
	u      *url.URL
}

// Implementation of apikey for DefaultClient
func (dc DefaultClientData) Apikey() (apikey string, err error) {
	if len(dc.apikey) < 1 {
		dc.apikey = os.Getenv(redashApikeyEnv)
	}
	if len(dc.apikey) < 1 {
		return "", errors.New("invalid apikey")
	}
	dc.Logger.Printf("[DEBUG] apikey: [%s]", maskKey(dc.apikey))
	return dc.apikey, nil
}

// Implementation of Url for DefaultClient
func (dc DefaultClientData) Url() (u *url.URL, err error) {
	if dc.u.String() == "" {
		dc.u, err = url.Parse(os.Getenv(redashUrlEnv))
		if err != nil {
			return nil, err
		}
	}
	return dc.u, err
}

// Implementation of HTTPClient for DefaultClient
func (dc DefaultClientData) HTTPClient() *http.Client {
	return http.DefaultClient
}

// Implementation of DefaultOpts for DefaultClient
func (dc DefaultClientData) DefaultOpts() *Options {
	return defaultOpts()
}

// Create a new defaultClient
func NewDefaultClient() *DefaultClientData {
	var u *url.URL
	var err error
	if ue := os.Getenv(redashUrlEnv); ue != "" {
		u, err = url.Parse(os.Getenv(redashUrlEnv))
		if err != nil {
			return nil
		}
		log.Printf("[DEBUG] set url: %s", u)
	} else {
		u = &url.URL{}
	}
	dcd := &DefaultClientData{
		apikey: os.Getenv(redashApikeyEnv),
		u:      u,
	}
	dcd.Logger = &log.Logger{}
	dcd.Logger.SetOutput(os.Stdout)
	dcd.Logger.SetFlags(log.Ldate | log.Ltime)
	return dcd
}

// maskKey is mask helper for apikey.
func maskKey(s string) string {
	var pre string
	if len(s) >= 4 {
		pre = s[0:4]
	} else {
		pre = "****"
	}
	return pre + "****"
}
