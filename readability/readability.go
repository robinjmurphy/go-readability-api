// Package readability provides access to the Readability APIs
package readability

import (
	"errors"
	"fmt"
	"github.com/garyburd/go-oauth/oauth"
	"net/http"
	"net/url"
)

const DefaultLoginURL = "https://www.readability.com/api/rest/v1/oauth/access_token/"
const DefaultReaderBaseURL = "https://www.readability.com/api/rest/v1/"

// A Client manages communication with the Readability APIs
type Client struct {
	LoginURL       string
	ReaderBaseURL  string
	ConsumerKey    string
	ConsumerSecret string
}

// A ReaderClient models the Readability Reader API
type ReaderClient struct {
	BaseURL          string
	OAuthClient      *oauth.Client
	OAuthCredentials *oauth.Credentials
}

// NewClient returns a new Readability client
func NewClient(key, secret string) *Client {
	client := Client{
		LoginURL:       DefaultLoginURL,
		ReaderBaseURL:  DefaultReaderBaseURL,
		ConsumerKey:    key,
		ConsumerSecret: secret,
	}
	return &client
}

// NewReader returns a new ReaderClient
func (client *Client) NewReaderClient(token, secret string) *ReaderClient {
	consumerCredentials := oauth.Credentials{Token: client.ConsumerKey, Secret: client.ConsumerSecret}
	userCredentials := oauth.Credentials{Token: token, Secret: secret}
	reader := ReaderClient{
		BaseURL:          client.ReaderBaseURL,
		OAuthClient:      &oauth.Client{Credentials: consumerCredentials},
		OAuthCredentials: &userCredentials,
	}
	return &reader
}

// Sign appends OAuth parameters to a set of url.Values
func (reader *ReaderClient) Sign(method string, uri string, data url.Values) {
	reader.OAuthClient.SignForm(reader.OAuthCredentials, "POST", uri, data)
}

// Post makes a HTTP POST request to the Reader API
func (reader *ReaderClient) Post(path string, data url.Values) (r *http.Response, err error) {
	uri := reader.BaseURL + path
	reader.Sign("POST", uri, data)
	resp, err := http.PostForm(uri, data)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode >= 400 {
		return resp, errors.New(fmt.Sprintf("%s. %s %s.", resp.Status, resp.Request.Method, resp.Request.URL))
	}
	return resp, nil
}

// AddBookmark creates a new bookmark for a URL
func (reader *ReaderClient) AddBookmark(uri string) (resp *http.Response, err error) {
	return reader.Post("/bookmarks", url.Values{"url": {uri}})
}
