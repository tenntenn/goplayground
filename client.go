package goplayground

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	// FrontBaseURL is frontend of the Go Playground.
	FrontBaseURL = "https://go.dev/play"
	// BaseURL is the default base URL of the Go Playground.
	BaseURL = "https://play.golang.org"
	// Deprecated: Go2GoBaseURL is the base URL of go2goplay.golang.org.
	Go2GoBaseURL = "https://go2goplay.golang.org"
	// Version is version of using Go Playground.
	Version = "2"
)

// Client is a client of Go Playground.
// If BaseURL is empty, Client uses default BaseURL.
// HTTPClient can be set instead of http.DefaultClient.
type Client struct {
	FrontBaseURL string
	BaseURL      string
	Backend      Backend
	HTTPClient   HTTPClient
}

func (cli *Client) baseURL() string {
	baseURL := BaseURL
	if cli.BaseURL != "" {
		baseURL = cli.BaseURL
	}

	if cli.Backend != BackendGotip {
		return baseURL
	}

	urlBase, err := url.Parse(baseURL)
	if err != nil {
		return baseURL
	}

	urlBase.Host = cli.Backend.String() + urlBase.Host

	return urlBase.String()
}

func (cli *Client) frontBaseURL() string {
	if cli.FrontBaseURL != "" {
		return cli.FrontBaseURL
	}
	return FrontBaseURL
}

func (cli *Client) httpClient() HTTPClient {
	if cli.HTTPClient != nil {
		return cli.HTTPClient
	}
	return http.DefaultClient
}

func srcToString(src interface{}) (string, error) {
	switch src := src.(type) {
	case io.Reader:
		bs, err := ioutil.ReadAll(src)
		if err != nil {
			return "", err
		}
		return srcToString(bs)
	case []byte:
		return string(src), nil
	case string:
		return src, nil
	}
	return "", errors.New("does not support src type")
}

func srcToReader(src interface{}) (io.Reader, error) {
	switch src := src.(type) {
	case io.Reader:
		return src, nil
	case []byte:
		return bytes.NewReader(src), nil
	case string:
		return strings.NewReader(src), nil
	}
	return nil, errors.New("does not support src type")
}
