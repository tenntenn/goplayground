package goplayground

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
)

// Download downloads source code hosted on Playground.
// The source would be written into w.
func (cli *Client) Download(w io.Writer, hashOrURL string) error {
	dlURL := cli.createDownloadURL(hashOrURL)
	req, err := http.NewRequest(http.MethodGet, dlURL, nil)
	if err != nil {
		return err
	}

	resp, err := cli.httpClient().Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("cannot download %s with %s", dlURL, resp.Status)
	}

	if _, err := io.Copy(w, resp.Body); err != nil {
		return err
	}

	return nil
}

func (cli *Client) createDownloadURL(hashOrURL string) string {
	switch {
	case strings.HasPrefix(hashOrURL, cli.baseURL()+"/p/"):
		return hashOrURL + ".go"
	case strings.HasPrefix(hashOrURL, cli.frontBaseURL()+"/p/"):
		dlURL, err := url.Parse(hashOrURL)
		if err == nil {
			hash := path.Base(dlURL.Path)
			if !strings.HasSuffix(hash, ".go") {
				hash += ".go"
			}
			return cli.baseURL() + "/p/" + hash
		}
	}
	// hash
	return cli.baseURL() + "/p/" + hashOrURL + ".go"
}
