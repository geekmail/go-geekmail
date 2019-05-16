package geekmail

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

type Conf struct {
	GitHubAuth GitHubAuth `yaml:"githubauth"`
	APIToken   string     `yaml:"apitoken"`
}

type Client struct {
	client *http.Client
	conf   *Conf
	// Base URL for API requests. baseURL should always be specified with a trailing slash.
	baseURL *url.URL

	// Services used for talking to different parts of the GeekMail API.
	Draft *DraftService
}

type service struct {
	client *Client
}

func NewClient(httpClient *http.Client, conf *Conf) *Client {
	baseURL, _ := url.Parse(defaultBaseURL)
	if httpClient == nil {
		httpClient = http.DefaultClient
	}
	c := &Client{
		client:  httpClient,
		conf:    conf,
		baseURL: baseURL,
	}

	s := &service{client: c}
	c.Draft = (*DraftService)(s)

	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, urlStr string, body interface{}) (*http.Request, error) {
	if !strings.HasSuffix(c.baseURL.Path, "/") {
		return nil, fmt.Errorf("baseURL must have a trailing slash, but %q does not", c.baseURL)
	}
	u, err := c.baseURL.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+c.conf.APIToken)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it. If rate limit is exceeded and reset time is in the future,
// Do returns *RateLimitError immediately without making a network API call.
//
// The provided ctx must be non-nil. If it is canceled or times out,
// ctx.Err() will be returned.
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// If the error type is *url.Error
		if e, ok := err.(*url.Error); ok {
			if url, err := url.Parse(e.URL); err == nil {
				e.URL = url.String()
				return nil, e
			}
		}

		return nil, err
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := checkError(resp, data); err != nil {
		return resp, err
	}

	if v != nil {
		if err := json.Unmarshal(data, v); err != nil {
			return resp, err
		}
	}

	return resp, nil
}

func checkError(r *http.Response, data []byte) error {
	reply := &APIResponse{}
	if data != nil {
		if err := json.Unmarshal(data, reply); err != nil {
			return err
		}
	}

	if reply.Code/100 != 2 {
		return fmt.Errorf(reply.Message)
	}

	return nil
}
