package goplayground

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	// BaseURL is the default base URL of Go Playground.
	BaseURL = "https://play.golang.org"
	// Go2BaseURL is the base URL for -go2 option.
	Go2BaseURL = "https://gotipplay.golang.org"
	// Version is version of using Go Playground.
	Version = "2"
)

// Client is a client of Go Playground.
// If BaseURL is empty, Client uses default BaseURL.
// HTTPClient can be set instead of http.DefaultClient.
type Client struct {
	BaseURL    string
	HTTPClient HTTPClient
}

func (cli *Client) baseURL() string {
	if cli.BaseURL != "" {
		return cli.BaseURL
	}
	return BaseURL
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
