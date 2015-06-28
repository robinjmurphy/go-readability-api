package readability

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	mux    *http.ServeMux
	client *Client
	reader *ReaderClient
	server *httptest.Server
)

func setup() {
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)
	client = NewClient("key", "secret")
	client.LoginURL = server.URL
	client.ReaderBaseURL = server.URL
	reader = client.NewReaderClient("token", "secret")
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

func TestNewClient(t *testing.T) {
	c := NewClient("key", "secret")
	readerBaseUrl := c.ReaderBaseURL
	if readerBaseUrl != DefaultReaderBaseURL {
		t.Errorf("NewClient ReaderBaseURL is %v, expected %v", readerBaseUrl, DefaultReaderBaseURL)
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
	token, secret, err := client.Login("username", "password")
	check(t, err)
	if token != expectedToken {
		t.Errorf("Token %s, expected %s", token, expectedToken)
	}
	if secret != expectedSecret {
		t.Errorf("Secret %s, expected %s", secret, expectedSecret)
	}
}

func TestAddBookmark(t *testing.T) {
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
