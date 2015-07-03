package readability

import (
	"net/http"
	"net/url"

	"github.com/garyburd/go-oauth/oauth"
)

// A ReaderClient models the Readability Reader API.
type ReaderClient struct {
	BaseURL          string
	OAuthClient      *oauth.Client
	OAuthCredentials *oauth.Credentials
}

// Post makes a HTTP POST request to the Reader API.
func (reader *ReaderClient) Post(path string, data url.Values) (r *http.Response, err error) {
	uri := reader.BaseURL + path
	return postWithCredentials(reader.OAuthClient, reader.OAuthCredentials, uri, data)
}

// AddBookmark creates a new bookmark for a URL.
func (reader *ReaderClient) AddBookmark(uri string) (r *http.Response, err error) {
	return reader.Post("/bookmarks", url.Values{"url": {uri}})
}
