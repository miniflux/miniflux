// Copyright 2017 Frédéric Guillot. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package http

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/miniflux/miniflux/errors"
	"github.com/miniflux/miniflux/logger"
	"github.com/miniflux/miniflux/timer"
	"github.com/miniflux/miniflux/version"
)

const (
	// 20 seconds max.
	requestTimeout = 20

	// 15MB max.
	maxBodySize = 1024 * 1024 * 15
)

var (
	errInvalidCertificate        = "Invalid SSL certificate (original error: %q)"
	errTemporaryNetworkOperation = "This website is temporarily unreachable (original error: %q)"
	errPermanentNetworkOperation = "This website is permanently unreachable (original error: %q)"
	errRequestTimeout            = "Website unreachable, the request timed out after %d seconds"
)

// Client is a HTTP Client :)
type Client struct {
	url                 string
	etagHeader          string
	lastModifiedHeader  string
	authorizationHeader string
	username            string
	password            string
	Insecure            bool
}

// Get execute a GET HTTP request.
func (c *Client) Get() (*Response, error) {
	request, err := c.buildRequest(http.MethodGet, nil)
	if err != nil {
		return nil, err
	}

	return c.executeRequest(request)
}

// PostForm execute a POST HTTP request with form values.
func (c *Client) PostForm(values url.Values) (*Response, error) {
	request, err := c.buildRequest(http.MethodPost, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	return c.executeRequest(request)
}

// PostJSON execute a POST HTTP request with JSON payload.
func (c *Client) PostJSON(data interface{}) (*Response, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	request, err := c.buildRequest(http.MethodPost, bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	request.Header.Add("Content-Type", "application/json")
	return c.executeRequest(request)
}

func (c *Client) executeRequest(request *http.Request) (*Response, error) {
	defer timer.ExecutionTime(time.Now(), fmt.Sprintf("[HttpClient] url=%s", c.url))

	client := c.buildClient()
	resp, err := client.Do(request)
	if err != nil {
		if uerr, ok := err.(*url.Error); ok {
			switch uerr.Err.(type) {
			case x509.CertificateInvalidError, x509.HostnameError:
				err = errors.NewLocalizedError(errInvalidCertificate, uerr.Err)
			case *net.OpError:
				if uerr.Err.(*net.OpError).Temporary() {
					err = errors.NewLocalizedError(errTemporaryNetworkOperation, uerr.Err)
				} else {
					err = errors.NewLocalizedError(errPermanentNetworkOperation, uerr.Err)
				}
			case net.Error:
				nerr := uerr.Err.(net.Error)
				if nerr.Timeout() {
					err = errors.NewLocalizedError(errRequestTimeout, requestTimeout)
				} else if nerr.Temporary() {
					err = errors.NewLocalizedError(errTemporaryNetworkOperation, nerr)
				}
			}
		}

		return nil, err
	}

	if resp.ContentLength > maxBodySize {
		return nil, fmt.Errorf("client: response too large (%d bytes)", resp.ContentLength)
	}

	response := &Response{
		Body:          resp.Body,
		StatusCode:    resp.StatusCode,
		EffectiveURL:  resp.Request.URL.String(),
		LastModified:  resp.Header.Get("Last-Modified"),
		ETag:          resp.Header.Get("ETag"),
		ContentType:   resp.Header.Get("Content-Type"),
		ContentLength: resp.ContentLength,
	}

	logger.Debug("[HttpClient:%s] OriginalURL=%s, StatusCode=%d, ContentLength=%d, ContentType=%s, ETag=%s, LastModified=%s, EffectiveURL=%s",
		request.Method,
		c.url,
		response.StatusCode,
		resp.ContentLength,
		response.ContentType,
		response.ETag,
		response.LastModified,
		response.EffectiveURL,
	)

	return response, err
}

func (c *Client) buildRequest(method string, body io.Reader) (*http.Request, error) {
	request, err := http.NewRequest(method, c.url, body)
	if err != nil {
		return nil, err
	}

	if c.username != "" && c.password != "" {
		request.SetBasicAuth(c.username, c.password)
	}

	request.Header = c.buildHeaders()
	return request, nil
}

func (c *Client) buildClient() http.Client {
	client := http.Client{Timeout: time.Duration(requestTimeout * time.Second)}
	if c.Insecure {
		client.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
	}

	return client
}

func (c *Client) buildHeaders() http.Header {
	headers := make(http.Header)
	headers.Add("User-Agent", "Mozilla/5.0 (compatible; Miniflux/"+version.Version+"; +https://miniflux.net)")
	headers.Add("Accept", "*/*")

	if c.etagHeader != "" {
		headers.Add("If-None-Match", c.etagHeader)
	}

	if c.lastModifiedHeader != "" {
		headers.Add("If-Modified-Since", c.lastModifiedHeader)
	}

	if c.authorizationHeader != "" {
		headers.Add("Authorization", c.authorizationHeader)
	}

	return headers
}

// NewClient returns a new HTTP client.
func NewClient(url string) *Client {
	return &Client{url: url, Insecure: false}
}

// NewClientWithCredentials returns a new HTTP client that requires authentication.
func NewClientWithCredentials(url, username, password string) *Client {
	return &Client{url: url, Insecure: false, username: username, password: password}
}

// NewClientWithAuthorization returns a new client with a custom authorization header.
func NewClientWithAuthorization(url, authorization string) *Client {
	return &Client{url: url, Insecure: false, authorizationHeader: authorization}
}

// NewClientWithCacheHeaders returns a new HTTP client that send cache headers.
func NewClientWithCacheHeaders(url, etagHeader, lastModifiedHeader string) *Client {
	return &Client{url: url, etagHeader: etagHeader, lastModifiedHeader: lastModifiedHeader, Insecure: false}
}
