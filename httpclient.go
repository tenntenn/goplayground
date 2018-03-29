package goplayground

import "net/http"

// HTTPClient is an interface of minimum HTTP client.
// net/http.Client implements this interface.
type HTTPClient interface {
	// Do method send a HTTP request.
	Do(*http.Request) (*http.Response, error)
}

// HTTPClientFunc implements HTTPClient.
type HTTPClientFunc func(*http.Request) (*http.Response, error)

// Do implements HTTPClient.Do.
func (f HTTPClientFunc) Do(req *http.Request) (*http.Response, error) {
	return f(req)
}
