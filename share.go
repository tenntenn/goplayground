package goplayground

import (
	"io/ioutil"
	"net/http"
	"net/url"
)

// Share generates share URL of the given src.
// src can be set string, []byte and io.Reader value.
func (cli *Client) Share(src interface{}) (*url.URL, error) {
	r, err := srcToReader(src)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, cli.baseURL()+"/share", r)
	resp, err := cli.httpClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	shareURL, err := url.Parse(cli.baseURL() + "/p/" + string(bs))
	if err != nil {
		return nil, err
	}

	return shareURL, nil
}
