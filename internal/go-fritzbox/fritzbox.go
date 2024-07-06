// Copyright 2015 The go-fritzbox AUTHORS. All rights reserved.
//
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package go_fritzbox

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultBaseURL = "http://fritz.box/"
)

// A Client manages communication with the FRITZ!Box
type Client struct {
	// HTTP client used to communicate with the FRITZ!Box
	client *http.Client

	// Base URL for requests. Defaults to the local fritzbox, but
	// can be set to a domain endpoint to use with an external FRITZ!Box.
	// BaseURL should always be specified with a trailing slash.
	BaseURL *url.URL

	// Session used to authenticate client
	Session *Session
}

// NewClient returns a new FRITZ!Box client. If a nil httpClient is
// provided, http.DefaultClient will be used. To use an external
// FRITZ!Box with a self-signed certificate, provide an http.Client
// that will be able to perform insecure connections (such as
// InsecureSkipVerify flag).
func NewClient(httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		client:  httpClient,
		BaseURL: baseURL,
	}

	return c
}

// NewRequest creates an API request. A relative URL can be provided
// in urlStr in which case it is resolved relative to the BaseURL of
// the Client. Relative URLs should always be specified without a
// preceding slash. If specified, the value pointed to by data is Query
// encoded and included as the request body in order to perform form requests.
func (c *Client) NewRequest(method, urlStr string, data url.Values) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	if c.Session != nil {
		values := u.Query()
		values.Set("sid", c.Session.Sid)
		u.RawQuery = values.Encode()
	}

	var buf io.Reader
	if data != nil {
		buf = strings.NewReader(data.Encode())
	}
	req, err := http.NewRequest(method, u.String(), buf)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return nil, err
	}
	return req, nil
}

// Do sends a request and returns the response. The response is
// either JSON decoded or XML encoded and stored in the value
// pointed to by v, or returned as an error, if any.
func (c *Client) Do(req *http.Request, v interface{}) (*http.Response, error) {
	if c.Session != nil {
		if err := c.Session.Refresh(); err != nil {
			return nil, err
		}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if c := resp.StatusCode; 200 < c && c > 299 {
		return nil, fmt.Errorf("wrong status code: %v", resp.StatusCode)
	}

	contentType := resp.Header.Get("Content-Type")

	if v != nil {
		if strings.Contains(contentType, "text/xml") {
			err = xml.NewDecoder(resp.Body).Decode(v)
		}
		if strings.Contains(contentType, "application/json") {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}

	return resp, err
}

// Auth sends a auth request and returns an error, if any. Session is stored
// in client in order to perform requests with authentification.
func (c *Client) Auth(username, password string) error {
	var s *Session
	if c.Session == nil {
		s = NewSession(c)
		c.Session = s
	} else {
		s = c.Session
	}

	if err := s.Open(); err != nil {
		return err
	}

	if err := s.Auth(username, password); err != nil {
		return err
	}

	return nil
}

func (c *Client) String() string {
	return c.Session.String()
}
