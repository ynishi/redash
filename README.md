# Redash(for Go)

Redash(for Go) is a unofficial rest api client lib.

## Current status

* Development.
* Some features can use.
* APIs maybe change(public function Get/Post/Delete won't change)

## Example

```
$ REDASH_APIKEY="..."
$ REDASH_URL="http://..."

package main

import "github.com/ynishi/redash"

client = redash.DefaultClient
byteJson = client.Get("/api/dashboards", nil) 
```

## Install 

```
$ go get "github.com/ynishi/redash"
```

## Development

```
$ git clone https://github.com/ynishi/redash.git 
$ cd redash
$ go test
```

## Credit and Licence

Copyright (c) 2017, Yutaka Nishimura.
Licenced under MIT, see LICENSE.
