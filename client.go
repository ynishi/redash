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

type Interface interface {
	Urler
	Apikeyer
	HTTPClienter
	DefaultOptser
}

// Get with Interface
func GetInter(data Interface, sub string, params map[string]string) (resp *http.Response, err error) {
	opts := data.DefaultOpts()
	for key, value := range params {
		opts.Params[key] = value
	}
	return DoInter(data, http.MethodGet, sub, opts)
}

// Post with Interface
func PostInter(data Interface, sub string, jsonBody []byte) (resp *http.Response, err error) {
	opts := data.DefaultOpts()
	for key, value := range defaultPostHeader {
		opts.Header[key] = value
	}
	opts.Body = bytes.NewReader(jsonBody)
	return DoInter(data, http.MethodPost, sub, opts)
}

// Delete with Interface
func DeleteInter(data Interface, sub string, params map[string]string) (resp *http.Response, err error) {
	opts := data.DefaultOpts()
	for key, value := range params {
		opts.Params[key] = value
	}
	return DoInter(data, http.MethodDelete, sub, opts)
}

// Do with Interface
func DoInter(data Interface, method, sub string, opts *Options) (resp *http.Response, err error) {
	req, err := RequestInter(data, method, sub, opts)
	if err != nil {
		return nil, err
	}
	return data.HTTPClient().Do(req)
}

// Request with Interface
func RequestInter(data Interface, method, sub string, opts *Options) (req *http.Request, err error) {
	u, err := data.Url()
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
	apikey, err := data.Apikey()
	if err != nil {
		return nil, err
	}
	opts.Header["Authorization"] = "Key " + apikey
	for key, value := range opts.Header {
		req.Header.Set(key, value)
	}
	return req, nil
}

// Get do Redash api GET
func Get(sub string, params map[string]string) (resp *http.Response, err error) {
	return GetInter(DefaultClient, sub, params)
}

// Post do Redash api POST
func Post(sub string, jsonBody []byte) (resp *http.Response, err error) {
	return PostInter(DefaultClient, sub, jsonBody)
}

// Delete do Redash api DELETE
func Delete(sub string, params map[string]string) (resp *http.Response, err error) {
	return DeleteInter(DefaultClient, sub, params)
}

// Do Redash api
func Do(method, sub string, opts *Options) (resp *http.Response, err error) {
	return DoInter(DefaultClient, method, sub, opts)
}

// Request make http.Request for Redash
func Request(method, sub string, opts *Options) (req *http.Request, err error) {
	return RequestInter(DefaultClient, method, sub, opts)
}

type ClientData struct {
	*log.Logger
}

type DefaultClientData struct {
	ClientData
	apikey string
	u      *url.URL
}

func (dc DefaultClientData) Apikey() (apikey string, err error) {
	if len(dc.apikey) < 1 {
		dc.apikey = os.Getenv(redashApikeyEnv)
	}
	if len(dc.apikey) < 1 {
		return "", errors.New("invalid apikey")
	}
	return dc.apikey, nil
}

func (dc DefaultClientData) Url() (u *url.URL, err error) {
	if dc.u.String() == "" {
		dc.u, err = url.Parse(os.Getenv(redashUrlEnv))
		if err != nil {
			return nil, err
		}
	}
	return dc.u, err
}

func (dc DefaultClientData) HTTPClient() *http.Client {
	return http.DefaultClient
}

func (dc DefaultClientData) DefaultOpts() *Options {
	return defaultOpts()
}

func NewDefaultClient() *DefaultClientData {
	var u *url.URL
	var err error
	if ue := os.Getenv(redashUrlEnv); ue != "" {
		u, err = url.Parse(os.Getenv(redashUrlEnv))
		if err != nil {
			return nil
		}
	} else {
		u = &url.URL{}
	}
	dcd := &DefaultClientData{
		apikey: os.Getenv(redashApikeyEnv),
		u:      u,
	}
	dcd.Logger = &log.Logger{}
	return dcd
}
