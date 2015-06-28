// Package readability provides access to the Readability APIs.
package readability

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/garyburd/go-oauth/oauth"
)

const DefaultLoginURL = "https://www.readability.com/api/rest/v1/oauth/access_token"
const DefaultReaderBaseURL = "https://www.readability.com/api/rest/v1"

// A Client manages communication with the Readability APIs.
type Client struct {
	LoginURL      string
	ReaderBaseURL string
	OAuthClient   *oauth.Client
}

// A ReaderClient models the Readability Reader API.
type ReaderClient struct {
	BaseURL          string
	OAuthClient      *oauth.Client
	OAuthCredentials *oauth.Credentials
}

// NewClient returns a new Readability client.
func NewClient(key, secret string) *Client {
	credentials := oauth.Credentials{Token: key, Secret: secret}
	client := Client{
		LoginURL:      DefaultLoginURL,
		ReaderBaseURL: DefaultReaderBaseURL,
		OAuthClient:   &oauth.Client{Credentials: credentials},
	}
	return &client
}

// NewReaderClient returns a new ReaderClient.
func (client *Client) NewReaderClient(token, secret string) *ReaderClient {
	credentials := oauth.Credentials{Token: token, Secret: secret}
	reader := ReaderClient{
		BaseURL:          client.ReaderBaseURL,
		OAuthClient:      client.OAuthClient,
		OAuthCredentials: &credentials,
	}
	return &reader
}

// Login returns an access token and secret for a user that can be used to
// create a ReaderClient.
func (client *Client) Login(username, password string) (token, secret string, err error) {
	resp, err := post(client.OAuthClient, nil, client.LoginURL, url.Values{
		"x_auth_username": {username},
		"x_auth_password": {password},
		"x_auth_mode":     {"client_auth"},
	})
	if err != nil {
		return token, secret, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if (err) != nil {
		return token, secret, err
	}
	formData, err := url.ParseQuery(string(body))
	if (err) != nil {
		return token, secret, err
	}
	return formData.Get("oauth_token"), formData.Get("oauth_token_secret"), nil
}

// Post makes a HTTP POST request to the Reader API.
func (reader *ReaderClient) Post(path string, data url.Values) (r *http.Response, err error) {
	uri := reader.BaseURL + path
	return post(reader.OAuthClient, reader.OAuthCredentials, uri, data)
}

// AddBookmark creates a new bookmark for a URL.
func (reader *ReaderClient) AddBookmark(uri string) (resp *http.Response, err error) {
	return reader.Post("/bookmarks", url.Values{"url": {uri}})
}

func post(client *oauth.Client, credentials *oauth.Credentials, uri string, data url.Values) (r *http.Response, err error) {
	client.SignForm(credentials, "POST", uri, data)
	resp, err := http.PostForm(uri, data)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode >= 400 {
		return resp, errors.New(
			fmt.Sprintf("Error %d %s: %s %s", resp.StatusCode, http.StatusText(resp.StatusCode), resp.Request.Method, resp.Request.URL))
	}
	return resp, nil
}
