package redash

import (
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

const (
	redashUrlEnv    = "REDASH_URL"
	redashApikeyEnv = "REDASH_APIKEY"
)

var (
	defaultOpts = Options{
		Params: make(map[string]string),
		Header: map[string]string{
			"User-Agent": "Go client",
		},
		Body: nil,
	}
	defaultPostHeader = map[string]string{
		"Content-Type":  "application/json",
	}
	DefaultClient = NewDefaultClient()
)

type Options struct {
	Params map[string]string
	Header map[string]string
	Body   io.Reader
}

type Interface interface {
	// Url is Redash server's endpoint.
	Url() (*url.URL, error)
	// Apikey is Redash Apikey to connect primary.
	Apikey() (string, error)
	// HTTPClient is HTTP client to do request.
	HTTPClient() *http.Client
	// DefaultOpts is default options for request.
	DefaultOpts() *Options
}

// Get with Interface
func GetInter(data Interface, sub string, params map[string]string) (resp *http.Response, err error) {
	opts := data.DefaultOpts()
	return DoInter(data, http.MethodGet, sub, *opts)
}

// Post with Interface
func PostInter(data Interface, sub string, params map[string][]string) (resp *http.Response, err error) {
	opts := data.DefaultOpts()
	for key, value := range defaultPostHeader {
		opts.Header[key] = value
	}
	return DoInter(data, http.MethodPost, sub, *opts)
}

// Delete with Interface
func DeleteInter(data Interface, sub string, params map[string][]string) (resp *http.Response, err error) {
	opts := data.DefaultOpts()
	return DoInter(data, http.MethodDelete, sub, *opts)
}

// Do with Interface
func DoInter(data Interface, method, sub string, opts Options) (resp *http.Response, err error) {
	req, err := RequestInter(data, method, sub, opts)
	if err != nil {
		return nil, err
	}
	return data.HTTPClient().Do(req)
}

// Request with Interface
func RequestInter(data Interface, method, sub string, opts Options) (req *http.Request, err error) {
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
func Post(sub string, params map[string][]string) (resp *http.Response, err error) {
	return PostInter(DefaultClient, sub, params)
}

// Delete do Redash api DELETE
func Delete(sub string, params map[string][]string) (resp *http.Response, err error) {
	return DeleteInter(DefaultClient, sub, params)
}

// Do Redash api
func Do(method, sub string, opts Options) (resp *http.Response, err error) {
	return DoInter(DefaultClient, method, sub, opts)
}

// Request make http.Request for Redash
func Request(method, sub string, opts Options) (req *http.Request, err error) {
	return RequestInter(DefaultClient, method, sub, opts)
}

type ClientData struct {
	Logger *log.Logger
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
	return &defaultOpts
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
