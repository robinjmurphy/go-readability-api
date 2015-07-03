package readability

import (
	"encoding/json"
	"io/ioutil"
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
