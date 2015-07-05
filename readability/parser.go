package readability

import (
	"net/http"
	"net/url"
)

// A ParserClient models the Readability Parser API.
type ParserClient struct {
	BaseURL string
	ApiKey  string
}

// Get makes a GET request to the Parser API
func (parser *ParserClient) Get(path string, query url.Values, v interface{}) (*http.Response, error) {
	query.Add("token", parser.ApiKey)
	uri := parser.BaseURL + path + "?" + query.Encode()
	return get(uri, v)
}

// Parse parses the contents of an article.
func (parser *ParserClient) Parse(articleURL string) (article Article, r *http.Response, err error) {
	query := url.Values{"url": {articleURL}}
	resp, err := parser.Get("/parser", query, &article)
	return article, resp, err
}

// Confidence returns the confidence with which an article can be parsed.
func (parser *ParserClient) Confidence(articleURL string) (confidence float64, r *http.Response, err error) {
	query := url.Values{"url": {articleURL}}
	parsed := Confidence{}
	resp, err := parser.Get("/confidence", query, &parsed)
	if err != nil {
		return confidence, resp, err
	}
	return parsed.Confidence, resp, nil
}
