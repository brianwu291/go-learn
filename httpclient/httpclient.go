package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// HTTPClient interface for mocking in tests
type (
	HTTPClient interface {
		Do(req *http.Request) (*http.Response, error)
	}

	Client struct {
		baseURL    string
		httpClient HTTPClient
	}

	Option func(*Client)

	Request struct {
		Method  string
		Path    string
		Query   map[string]string
		Headers map[string]string
		Body    io.Reader
	}
)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient HTTPClient) Option {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithBaseURL sets the base URL for the client
func WithBaseURL(baseURL string) Option {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// NewClient creates a new client with default configuration
// 30 secs timeout
func NewClient(opts ...Option) *Client {
	c := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func (c *Client) Get(ctx context.Context, path string, result interface{}) error {
	return c.Do(ctx, Request{
		Method: http.MethodGet,
		Path:   path,
	}, result)
}

func (c *Client) Post(ctx context.Context, path string, body io.Reader, result interface{}) error {
	return c.Do(ctx, Request{
		Method: http.MethodPost,
		Path:   path,
		Body:   body,
	}, result)
}

// Do executes an HTTP request and decodes the response
func (c *Client) Do(ctx context.Context, r Request, result interface{}) error {
	url := fmt.Sprintf("%s%s", c.baseURL, r.Path)

	req, err := http.NewRequestWithContext(ctx, r.Method, url, r.Body)
	if err != nil {
		return fmt.Errorf("creating request: %w", err)
	}

	// Add query parameters
	if len(r.Query) > 0 {
		q := req.URL.Query()
		for key, value := range r.Query {
			q.Add(key, value)
		}
		req.URL.RawQuery = q.Encode()
	}

	// Add headers
	for key, value := range r.Headers {
		req.Header.Add(key, value)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("decoding response: %w", err)
		}
	}

	return nil
}
