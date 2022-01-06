package goplayground

import (
	"encoding/json"
	"net/http"
)

// VersionResult is result of Client.Format.
type VersionResult struct {
	Version string
	Release string
	Name    string
}

// Version gets version and release tags which is used in the Go Playground.
func (cli *Client) Version() (*VersionResult, error) {
	req, err := http.NewRequest(http.MethodGet, cli.baseURL()+"/version", nil)
	if err != nil {
		return nil, err
	}

	resp, err := cli.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result VersionResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
