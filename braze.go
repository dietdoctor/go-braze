package braze

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	defaultBaseURL   = "https://rest.iad-05.braze.com"
	defaultUserAgent = "go-braze"
)

// Braze defines the Braze REST API client interface.
type Braze interface {
	UsersEndpoint
}

// Client implements Braze REST API client.
type Client struct {
	baseURL    *url.URL
	apiKey     string
	userAgent  string
	httpClient *http.Client

	// TODO
	// Export ExportService
	// Email EmailService
	// Subscription SubscriptionService
	// Templates    TemplatesService

	Messaging MessagingEndpoint
	Users     UsersEndpoint
}

// NewClient sets up a new Trustpilot client.
func NewClient(opts ...ClientOption) (*Client, error) {
	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: time.Second * 10,
		},
	}
	if err := c.applyOptions(opts...); err != nil {
		return nil, err
	}

	c.Users = &UsersService{client: c}
	c.Messaging = &MessagingService{client: c}

	return c, nil
}

// ClientOption is a functional option for configuring the API client.
type ClientOption func(*Client) error

// BaseURL allows to change the default API base url.
func BaseURL(u *url.URL) ClientOption {
	return func(c *Client) error {
		c.baseURL = u
		return nil
	}
}

// APIKey is a functional option for configuring api access key.
func APIKey(k string) ClientOption {
	return func(c *Client) error {
		c.apiKey = k
		return nil
	}
}

// UserAgent is a functional option for configuring client user agent.
func UserAgent(a string) ClientOption {
	return func(c *Client) error {
		c.userAgent = a
		return nil
	}
}

// HTTPClient is a functional option for configuring http client.
func HTTPClient(h *http.Client) ClientOption {
	return func(c *Client) error {
		c.httpClient = h
		return nil
	}
}

func (c *Client) applyOptions(opts ...ClientOption) error {
	for _, o := range opts {
		if err := o(c); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) newRequest(method string, path string, body interface{}) (*http.Request, error) {
	u := c.baseURL.ResolveReference(&url.URL{Path: path})

	var b []byte
	if body != nil {
		jb, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		b = jb
	}

	req, err := http.NewRequest(method, u.String(), bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Authorization", "Bearer "+c.apiKey)

	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request, v interface{}) error {
	// TODO remove
	// reqDump, err := httputil.DumpRequest(req, true)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%s\n", reqDump)

	resp, err := c.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TODO remove
	// respDump, err := httputil.DumpResponse(resp, true)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Printf("%s\n", respDump)

	if err := parseError(resp); err != nil {
		return err
	}

	if resp.ContentLength != 0 && v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return err
		}
	}

	return nil
}

// Only errors will be parsed into ErrorResponse.
func parseError(resp *http.Response) error {
	switch resp.StatusCode {
	case http.StatusOK, http.StatusCreated, http.StatusAccepted:
		return nil
	case http.StatusNotFound,
		http.StatusBadRequest,
		http.StatusUnprocessableEntity,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusTooManyRequests:
		e := &ErrorResponse{}
		if err := json.NewDecoder(resp.Body).Decode(&e); err != nil {
			return err
		}
		e.ErrorCode = resp.StatusCode
		return e
	default:
		// Don't assume every other error would have a valid json response object.
		return &ErrorResponse{ErrorCode: resp.StatusCode}
	}
}

func (r *ErrorResponse) Error() string {
	b := strings.Builder{}
	b.WriteString(fmt.Sprintf("%d: ", r.ErrorCode))
	b.WriteString(r.Message)
	if len(r.Errors) != 0 {
		b.WriteString(": ")
		b.WriteString(fmt.Sprintf("%v", r.Errors))
	}
	return b.String()
}

// ErrorResponse includes an ErrorCode as well.
type ErrorResponse struct {
	Response
	ErrorCode int
}

type Response struct {
	Message string  `json:"message,omitempty"`
	SendID  string  `json:"send_id,omitempty"`
	Deleted int     `json:"deleted,omitempty"`
	Errors  []Error `json:"errors,omitempty"` // Minor errors.
}

type Error struct {
	Type       string `json:"type,omitempty"`
	InputArray string `json:"input_array,omitempty"`
	Index      int    `json:"index,omitempty"`
}

// Bool is a helper routine that allocates a new bool value
// to store v and returns a pointer to it.
func Bool(v bool) *bool { return &v }

// Int is a helper routine that allocates a new int value
// to store v and returns a pointer to it.
func Int(v int) *int { return &v }

// Int64 is a helper routine that allocates a new int64 value
// to store v and returns a pointer to it.
func Int64(v int64) *int64 { return &v }

// String is a helper routine that allocates a new string value
// to store v and returns a pointer to it.
func String(v string) *string { return &v }

// Float64 is a helper routine that allocates a new float64 value
// to store v and returns a pointer to it.
func Float64(v float64) *float64 { return &v }
