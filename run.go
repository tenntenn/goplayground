package goplayground

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"
)

// RunResult is result of Client.Run.
type RunResult struct {
	// Errors is compile or runtime error on Go Playground.
	Errors string
	// Events has output events on Go Playground.
	Events []*RunEvent
}

// RunEvent represents output events to stdout or stderr of Client.Run.
type RunEvent struct {
	// Message is a message which is outputed to stdout or stderr.
	Message string
	// Kind has stdout or stderr value.
	Kind string
	// Delay represents delay time to print the message to stdout or stderr.
	Delay time.Duration
}

// Run compiles and runs the given src.
// src can be set string, []byte and io.Reader value.
func (cli *Client) Run(src interface{}) (*RunResult, error) {
	values := url.Values{}
	values.Set("version", Version)
	body, err := srcToString(src)
	if err != nil {
		return nil, err
	}
	values.Set("body", body)

	req, err := http.NewRequest(http.MethodPost, cli.baseURL()+"/compile?"+values.Encode(), nil)
	if err != nil {
		return nil, err
	}

	resp, err := cli.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var result RunResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result, nil
}
