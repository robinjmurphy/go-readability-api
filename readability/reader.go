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

// AddBookmark creates a new bookmark for a URL.
func (reader *ReaderClient) AddBookmark(uri string) (*http.Response, error) {
	return reader.Post("/bookmarks", url.Values{"url": {uri}}, nil)
}

// GetArticle retrieves a single article and its content.
func (reader *ReaderClient) GetArticle(articleId string) (article Article, resp *http.Response, err error) {
	path := "/articles/" + articleId
	resp, err = reader.Get(path, url.Values{}, &article)
	return article, resp, err
}

// Get makes a HTTP GET request to the Reader API.
func (reader *ReaderClient) Get(path string, query url.Values, v interface{}) (*http.Response, error) {
	uri := reader.BaseURL + path
	reader.OAuthClient.SignParam(reader.OAuthCredentials, "GET", uri, query)
	uri = uri + "?" + query.Encode()
	return get(uri, v)
}

// Post makes a HTTP POST request to the Reader API.
func (reader *ReaderClient) Post(path string, data url.Values, v interface{}) (*http.Response, error) {
	uri := reader.BaseURL + path
	reader.OAuthClient.SignForm(reader.OAuthCredentials, "POST", uri, data)
	return post(uri, data, v)
}
