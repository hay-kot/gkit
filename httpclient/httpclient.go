// Package httpclient provides a client for working with HTTP requests, using
// a method based API.
package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path"
	"strings"
)

type ClientMiddleware = func(*http.Request) (*http.Request, error)

type Client struct {
	client  *http.Client
	baseURL string
	mw      []ClientMiddleware
}

func New(client *http.Client, base string) *Client {
	return &Client{
		client:  client,
		baseURL: base,
		mw:      nil,
	}
}

func (c *Client) Use(mws ...ClientMiddleware) {
	c.mw = append(c.mw, mws...)
}

func (c *Client) Get(url string, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) GetCtx(ctx context.Context, url string, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) Post(url string, payload io.Reader, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) PostCtx(ctx context.Context, url string, payload io.Reader, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, payload)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) Put(url string, payload io.Reader, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, payload)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) PutCtx(ctx context.Context, url string, payload io.Reader, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, payload)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) Delete(url string, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

func (c *Client) DeleteCtx(ctx context.Context, url string, mw ...ClientMiddleware) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	return c.Do(req, mw)
}

// Do will execute the provided request, applying any middleware that has been
// registered with the client. Any middleware provided to this function will
// be applied after the middleware registered with the client.
//
// Example:
//
//	Request -> ClientMiddleware(Request) -> FunctionMiddleware(ClientMiddleware(Request))
func (c *Client) Do(req *http.Request, rmw []ClientMiddleware) (*http.Response, error) {
	for _, mw := range c.mw {
		var err error
		req, err = mw(req)
		if err != nil {
			return nil, err
		}
	}

	for _, mw := range rmw {
		var err error
		req, err = mw(req)
		if err != nil {
			return nil, err
		}
	}

	return c.client.Do(req)
}

// Path will safely join the base URL and the provided path and return a string
// that can be used in a request.
func (c *Client) Path(url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	}

	base := strings.TrimRight(c.baseURL, "/")
	if url == "" {
		return base
	}

	return path.Join(base, strings.TrimLeft(url, "/"))
}

// Pathf will call fmt.Sprintf with the provided values and then pass them
// to Client.Path as a convenience.
func (c *Client) Pathf(url string, v ...any) string {
	url = fmt.Sprintf(url, v...)
	return c.Path(url)
}

// DecodeJSON will decode the response body into the provided value.
func DecodeJSON[T any](r *http.Response) (T, error) {
	var zero T

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&zero); err != nil {
		return zero, err
	}
	return zero, nil
}
