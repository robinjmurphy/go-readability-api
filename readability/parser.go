package readability

import (
	"encoding/json"
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
	body, resp, err := get(parser.BaseURL+"/parser", url.Values{
		"url":   {articleURL},
		"token": {parser.ApiKey},
	})
	if err != nil {
		return article, resp, err
	}
	err = json.Unmarshal(body, &article)
	if err != nil {
		return article, resp, err
	}
	return article, resp, nil
}

// Confidence returns the confidence with which an article can be parsed.
func (parser *ParserClient) Confidence(articleURL string) (confidence float64, r *http.Response, err error) {
	body, resp, err := get(parser.BaseURL+"/confidence", url.Values{
		"url":   {articleURL},
		"token": {parser.ApiKey},
	})
	if err != nil {
		return confidence, resp, err
	}
	parsed := Confidence{}
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		return confidence, resp, err
	}
	return parsed.Confidence, resp, nil
}
