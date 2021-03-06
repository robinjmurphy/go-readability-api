package readability

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux    *http.ServeMux
	reader *ReaderClient
	parser *ParserClient
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	LoginURL = server.URL
	reader = NewReaderClient("consumer_token", "consumer_secret", "token", "secret")
	parser = NewParserClient("parser_api_key")
	reader.BaseURL = server.URL
	parser.BaseURL = server.URL
}

func teardown() {
	server.Close()
}

func testMethod(t *testing.T, r *http.Request, expected string) {
	if r.Method != expected {
		t.Errorf("Request method: %v, expected %v", r.Method, expected)
	}
}

func check(t *testing.T, err error) {
	if err != nil {
		t.Error(err)
	}
}

func TestLogin(t *testing.T) {
	setup()
	defer teardown()
	expectedToken := "a_token"
	expectedSecret := "a_secret"
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		fmt.Fprintf(w, "oauth_token=%s&oauth_token_secret=%s", expectedToken, expectedSecret)
	})
	token, secret, err := Login("consumer_token", "consumer_secret", "username", "password")
	check(t, err)
	if token != expectedToken {
		t.Errorf("Token %s, expected %s", token, expectedToken)
	}
	if secret != expectedSecret {
		t.Errorf("Secret %s, expected %s", secret, expectedSecret)
	}
}

func TestReaderAddBookmark(t *testing.T) {
	setup()
	defer teardown()
	expectedLocation := "https://www.readability.com/api/rest/v1/bookmarks/1"
	mux.HandleFunc("/bookmarks", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		w.Header().Set("Location", expectedLocation)
	})
	resp, err := reader.AddBookmark("http://www.example.com/")
	check(t, err)
	location := resp.Header.Get("location")
	if location != expectedLocation {
		t.Errorf("Location %v, expected %v", location, expectedLocation)
	}
}

func TestReaderGetArticle(t *testing.T) {
	setup()
	defer teardown()
	expectedAuthor := "Steve Jobs"
	expectedShortURL := "http://rdd.me/4ksnrhhl"
	mux.HandleFunc("/articles/123", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"author": "%s", "short_url": "%s"}`, expectedAuthor, expectedShortURL)
	})
	article, _, err := reader.GetArticle("123")
	check(t, err)
	if article.Author != expectedAuthor {
		t.Errorf("Author %v, expected %v", article.Author, expectedAuthor)
	}
	if article.ShortURL != expectedShortURL {
		t.Errorf("ShortUrl %v, expected %v", article.ShortURL, expectedShortURL)
	}
}

func TestParserParse(t *testing.T) {
	setup()
	defer teardown()
	expectedAuthor := "Steve Jobs"
	expectedShortURL := "http://rdd.me/4ksnrhhl"
	mux.HandleFunc("/parser", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"author": "%s", "short_url": "%s"}`, expectedAuthor, expectedShortURL)
	})
	article, _, err := parser.Parse("http://www.example.com/")
	check(t, err)
	if article.Author != expectedAuthor {
		t.Errorf("Author %v, expected %v", article.Author, expectedAuthor)
	}
	if article.ShortURL != expectedShortURL {
		t.Errorf("ShortUrl %v, expected %v", article.ShortURL, expectedShortURL)
	}
}

func TestParserConfidence(t *testing.T) {
	setup()
	defer teardown()
	expectedConfidence := 5.5
	mux.HandleFunc("/confidence", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprintf(w, `{"url": "http://www.example.com/", "confidence": %v}`, expectedConfidence)
	})
	confidence, _, err := parser.Confidence("http://www.example.com/")
	check(t, err)
	if confidence != expectedConfidence {
		t.Errorf("Confidence %v, expected %v", confidence, expectedConfidence)
	}
}

func TestParserConfidence_invalidJson(t *testing.T) {
	setup()
	defer teardown()
	mux.HandleFunc("/confidence", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		fmt.Fprint(w, `foo`)
	})
	_, _, err := parser.Confidence("http://www.example.com/")
	if err == nil {
		t.Error("Expected error")
	}
}
