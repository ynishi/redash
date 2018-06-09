# Redash(for Go)

Redash(for Go) is a unofficial simple api client lib.

See Godoc at http://godoc.org/github.com/ynishi/redash

Redash is OSS BI tool. See more at https://redash.io/

## Current status

* Version 1.0(v1)

## Example

### set env

```shell
$ export REDASH_APIKEY="abc..."
$ export REDASH_URL="http://localhost"
```

### code

```go
package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/ynishi/redash/v1"
)

func main() {
	response, err := redash.Get("/api/queries", nil)
	if err != nil {
		log.Fatal(err)
	}
	buf, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%v", string(buf))
}
```

## Install

```shell
$ go get "github.com/ynishi/redash"
```

## Development

Welcome to participate develop, send pull request, add issue(question, bugs, wants and so on).

### Start develop

* first, clone repository.
```shell
$ git clone https://github.com/ynishi/redash.git
$ cd redash
$ go test
```
* and make pull request.

## Credit and License

Copyright (c) 2017, Yutaka Nishimura.
Licensed under MIT, see LICENSE.
