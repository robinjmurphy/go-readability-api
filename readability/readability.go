// Package readability provides access to the Readability APIs.
package readability

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/garyburd/go-oauth/oauth"
)

const DefaultLoginURL = "https://www.readability.com/api/rest/v1/oauth/access_token"
const DefaultReaderBaseURL = "https://www.readability.com/api/rest/v1"
const DefaultParserBaseURL = "https://www.readability.com/api/content/v1"

// A Client manages communication with the Readability APIs.
type Client struct {
	LoginURL      string
	ReaderBaseURL string
	ParserBaseURL string
	ParserApiKey  string
	OAuthClient   *oauth.Client
}

// NewClient returns a new Readability client.
func NewClient(key, secret, parserApiKey string) *Client {
	credentials := oauth.Credentials{Token: key, Secret: secret}
	client := Client{
		LoginURL:      DefaultLoginURL,
		ReaderBaseURL: DefaultReaderBaseURL,
		ParserBaseURL: DefaultParserBaseURL,
		OAuthClient:   &oauth.Client{Credentials: credentials},
		ParserApiKey:  parserApiKey,
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

// NewParserClient returns a new ParserClient.
func (client *Client) NewParserClient() *ParserClient {
	parser := ParserClient{
		BaseURL: client.ParserBaseURL,
		ApiKey:  client.ParserApiKey,
	}
	return &parser
}

// Login returns an access token and secret for a user that can be used to
// create a ReaderClient.
func (client *Client) Login(username, password string) (token, secret string, err error) {
	data := url.Values{
		"x_auth_username": {username},
		"x_auth_password": {password},
		"x_auth_mode":     {"client_auth"},
	}
	client.OAuthClient.SignForm(nil, "POST", client.LoginURL, data)
	resp, err := post(client.LoginURL, data, nil)
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

func get(uri string, v interface{}) (*http.Response, error) {
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	return do(req, v)
}

func post(uri string, data url.Values, v interface{}) (*http.Response, error) {
	req, err := http.NewRequest("POST", uri, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	return do(req, v)
}

func do(req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode >= 400 {
		return resp, httpError(resp)
	}
	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
		if err != nil {
			return resp, err
		}
	}
	return resp, nil
}

func httpError(resp *http.Response) error {
	return errors.New(
		fmt.Sprintf("Error %d %s: %s %s", resp.StatusCode, http.StatusText(resp.StatusCode), resp.Request.Method, resp.Request.URL))
}
