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

// Parse parses the contents of an article.
func (parser *ParserClient) Parse(articleURL string) (article Article, r *http.Response, err error) {
	resp, err := get(parser.BaseURL+"/parser", url.Values{
		"url":   {articleURL},
		"token": {parser.ApiKey},
	}, &article)
	return article, resp, err
}

// Confidence returns the confidence with which an article can be parsed.
func (parser *ParserClient) Confidence(articleURL string) (confidence float64, r *http.Response, err error) {
	parsed := Confidence{}
	resp, err := get(parser.BaseURL+"/confidence", url.Values{
		"url":   {articleURL},
		"token": {parser.ApiKey},
	}, &parsed)
	if err != nil {
		return confidence, resp, err
	}
	return parsed.Confidence, resp, nil
}
