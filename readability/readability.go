// Package readability provides access to the Readability APIs.
package readability

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

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

// A ReaderClient models the Readability Reader API.
type ReaderClient struct {
	BaseURL          string
	OAuthClient      *oauth.Client
	OAuthCredentials *oauth.Credentials
}

// A ParserClient models the Readability Parser API.
type ParserClient struct {
	BaseURL string
	ApiKey  string
}

// An Article represents a Readability article object.
type Article struct {
	Author        interface{} `json:"author"`
	Content       string      `json:"content"`
	DatePublished interface{} `json:"date_published"`
	Dek           interface{} `json:"dek"`
	Direction     string      `json:"direction"`
	Domain        string      `json:"domain"`
	Excerpt       string      `json:"excerpt"`
	LeadImageURL  string      `json:"lead_image_url"`
	NextPageID    interface{} `json:"next_page_id"`
	RenderedPages int         `json:"rendered_pages"`
	ShortURL      string      `json:"short_url"`
	Title         string      `json:"title"`
	TotalPages    int         `json:"total_pages"`
	URL           string      `json:"url"`
	WordCount     int         `json:"word_count"`
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
func (reader *ReaderClient) AddBookmark(uri string) (r *http.Response, err error) {
	return reader.Post("/bookmarks", url.Values{"url": {uri}})
}

// Parse parses the contents of an article.
func (parser *ParserClient) Parse(articleURL string) (article Article, r *http.Response, err error) {
	resp, err := get(parser.BaseURL+"/parser", url.Values{
		"url":   {articleURL},
		"token": {parser.ApiKey},
	})
	if err != nil {
		return article, resp, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return article, resp, err
	}
	err = json.Unmarshal(body, &article)
	if err != nil {
		return article, resp, err
	}
	return article, resp, nil
}

func post(client *oauth.Client, credentials *oauth.Credentials, uri string, data url.Values) (r *http.Response, err error) {
	client.SignForm(credentials, "POST", uri, data)
	resp, err := http.PostForm(uri, data)
	if err != nil {
		return resp, err
	}
	if resp.StatusCode >= 400 {
		return resp, httpError(resp)
	}
	return resp, nil
}

func get(uri string, query url.Values) (r *http.Response, err error) {
	resp, err := http.Get(uri + "?" + query.Encode())
	if err != nil {
		return r, err
	}
	if resp.StatusCode >= 400 {
		return resp, httpError(resp)
	}
	return resp, nil
}

func httpError(resp *http.Response) error {
	return errors.New(
		fmt.Sprintf("Error %d %s: %s %s", resp.StatusCode, http.StatusText(resp.StatusCode), resp.Request.Method, resp.Request.URL))
}
