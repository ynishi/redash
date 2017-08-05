# Redash(for Go)

Redash(for Go) is a unofficial simple api client lib.

See Godoc at http://godoc.org/github.com/ynishi/redash

Redash is OSS BI tool. See more at https://redash.io/

## Current status

* Version 1.0(v1)

## Example

### set env

```
$ export REDASH_APIKEY="abc..."
$ export REDASH_URL="http://localhost"
```

### code 

```
package main

redash "github.com/ynishi/redash/v1"

response, _ := redash.Get("/api/queries", nil)
buf := ioutil.ReadAll(response.Body)
fmt.Printf("%v", string(buf))
```

## Install 

```
$ go get "github.com/ynishi/redash"
```

## Development

Welcome to participate develop, send pull request, add issue(question, bugs, wants and so on).

### Start develop

* first, clone repository.
```
$ git clone https://github.com/ynishi/redash.git 
$ cd redash
$ go test
```
* and make rull request.

## Credit and License

Copyright (c) 2017, Yutaka Nishimura.
Licenced under MIT, see LICENSE.
