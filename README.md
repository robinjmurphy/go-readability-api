# go-readability-api

[![GoDoc](https://godoc.org/github.com/robinjmurphy/go-readability-api/readability?status.svg)](https://godoc.org/github.com/robinjmurphy/go-readability-api/readability)

> Go client library for accessing the [Readability APIs](https://readability.com/developers/api)

This package is a work in progress and currently **only** supports retrieving user credentials and creating bookmarks.

## Installation

```
go get github.com/robinjmurphy/go-readability-api/readability
```

## Usage

```go
package main

import "github.com/robinjmurphy/go-readability-api/readability"

func main() {
  client := readability.NewClient("consumer_key", "consumer_secret")
  reader := client.NewReaderClient("access_token", "access_token_secret")
  _, err := reader.AddBookmark("http://www.bbc.co.uk/news/technology-33228149")
  if err != nil {
    panic(err)
  }
}
```

### Authentication

You can retrieve an access token and secret for a user with their username and password:

```go
token, secret, _, err := client.Login("username", "password")
if err != nill {
  panic(err)
}

```

See the [full package documentation](https://godoc.org/github.com/robinjmurphy/go-readability-api/readability) for the complete API.
