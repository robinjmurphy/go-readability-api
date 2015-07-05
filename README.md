# go-readability-api

[![GoDoc](https://godoc.org/github.com/robinjmurphy/go-readability-api/readability?status.svg)](https://godoc.org/github.com/robinjmurphy/go-readability-api/readability) [![Build Status](https://travis-ci.org/robinjmurphy/go-readability-api.svg)](https://travis-ci.org/robinjmurphy/go-readability-api)

> Go client library for accessing the [Readability APIs](https://readability.com/developers/api)

This package is a work in progress and currently **only** supports retrieving user credentials, creating bookmarks and parsing articles.

## Installation

```
go get github.com/robinjmurphy/go-readability-api/readability
```

## Usage

### Using the Reader API

```go
package main

import "github.com/robinjmurphy/go-readability-api/readability"

func main() {
  reader := readability.NewReaderClient("consumer_key", "consumer_secret", "access_token", "access_token_secret")
  _, err := reader.AddBookmark("http://www.bbc.co.uk/news/technology-33228149")
  if err != nil {
    panic(err)
  }
}
```

### Using the Parser API

```go
package main

import (
  "fmt"

  "github.com/robinjmurphy/go-readability-api/readability"
)

func main() {
  parser := readability.NewParserClient("parser_api_key")
  article, _, err := parser.Parse("http://www.bbc.co.uk/news/technology-33228149")
  if err != nil {
    panic(err)
  }
  fmt.Println(article.Title)
}
```

### Authentication

You can retrieve an access token and secret for a user with their username and password:

```go
token, secret, _, err := readability.Login("consumer_key", "consumer_secret", "username", "password")
if err != nil {
  panic(err)
}

```

See the [full package documentation](https://godoc.org/github.com/robinjmurphy/go-readability-api/readability) for the complete API.
