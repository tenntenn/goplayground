package goplayground

import (
	"encoding/json"
	"net/http"
	"net/url"
)

// FormatResult is result of Client.Format.
type FormatResult struct {
	// Body is the formatted source code.
	Body string
	// Error is a gofmt error.
	Error string
}

// Format formats the given src by gofmt or goimports.
// src can be set string, []byte and io.Reader value.
// If imports is true, Format formats and imports unimport packages with goimports.
func (cli *Client) Format(src interface{}, imports bool) (*FormatResult, error) {
	values := url.Values{}
	if imports {
		values.Set("imports", "true")
	}
	body, err := srcToString(src)
	if err != nil {
		return nil, err
	}
	values.Set("body", body)

	req, err := http.NewRequest(http.MethodPost, cli.baseURL()+"/fmt?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := cli.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result FormatResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
