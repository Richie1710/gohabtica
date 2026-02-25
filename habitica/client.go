package habitica

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/danielrichardt/gohabitica/internal/config"
)

// DefaultUserAgent is the user agent used when none is provided.
const DefaultUserAgent = "gohabitica-client/0.1"

// Client wraps access to the Habitica API.
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client

	userID   string
	apiToken string

	userAgent string
	clientID  string

	// Services
	User      *UserService
	Tasks     *TasksService
	Groups    *GroupsService
	Challenges *ChallengesService
	Content   *ContentService
	Tags      *TagsService
	Shops     *ShopsService
	Webhooks  *WebhooksService
	Admin     *AdminService
}

// Option configures the client during construction.
type Option func(*Client)

// WithHTTPClient overrides the HTTP client used by the Habitica client.
func WithHTTPClient(hc *http.Client) Option {
	return func(c *Client) {
		if hc != nil {
			c.httpClient = hc
		}
	}
}

// WithBaseURL overrides the base URL of the API.
func WithBaseURL(raw string) Option {
	return func(c *Client) {
		if raw == "" {
			return
		}
		u, err := url.Parse(raw)
		if err == nil {
			c.baseURL = u
		}
	}
}

// WithUserAgent sets a custom User-Agent header.
func WithUserAgent(ua string) Option {
	return func(c *Client) {
		if ua != "" {
			c.userAgent = ua
		}
	}
}

// WithClientID sets the x-client header (for identifying custom tools).
func WithClientID(id string) Option {
	return func(c *Client) {
		c.clientID = id
	}
}

// NewClient creates a new Habitica client based on the given configuration.
func NewClient(cfg *config.Config, opts ...Option) (*Client, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config darf nicht nil sein")
	}
	if cfg.BaseURL == "" {
		return nil, fmt.Errorf("config.BaseURL darf nicht leer sein")
	}
	if cfg.UserID == "" || cfg.APIToken == "" {
		return nil, fmt.Errorf("UserID und APIToken müssen gesetzt sein")
	}

	baseURL, err := url.Parse(cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("ungültige BaseURL %q: %w", cfg.BaseURL, err)
	}

	c := &Client{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		userID:    cfg.UserID,
		apiToken:  cfg.APIToken,
		userAgent: DefaultUserAgent,
		clientID:  "gohabitica-client",
	}

	for _, opt := range opts {
		opt(c)
	}

	// Initialize services
	c.User = &UserService{client: c}
	c.Tasks = &TasksService{client: c}
	c.Groups = &GroupsService{client: c}
	c.Challenges = &ChallengesService{client: c}
	c.Content = &ContentService{client: c}
	c.Tags = &TagsService{client: c}
	c.Shops = &ShopsService{client: c}
	c.Webhooks = &WebhooksService{client: c}
	c.Admin = &AdminService{client: c}

	return c, nil
}

// newRequest builds a new HTTP request relative to the base URL.
func (c *Client) newRequest(ctx context.Context, method, p string, query url.Values, body any) (*http.Request, error) {
	if ctx == nil {
		ctx = context.Background()
	}

	rel := &url.URL{Path: path.Join(c.baseURL.Path, p)}
	if query != nil {
		rel.RawQuery = query.Encode()
	}
	u := c.baseURL.ResolveReference(rel)

	var buf io.ReadWriter
	if body != nil {
		buf = &bytes.Buffer{}
		enc := json.NewEncoder(buf)
		if err := enc.Encode(body); err != nil {
			return nil, fmt.Errorf("body konnte nicht serialisiert werden: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("x-api-user", c.userID)
	req.Header.Set("x-api-key", c.apiToken)
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.userAgent != "" {
		req.Header.Set("User-Agent", c.userAgent)
	}
	if c.clientID != "" {
		req.Header.Set("x-client", c.clientID)
	}

	return req, nil
}

// doRequest executes an HTTP request and decodes the standardized Habitica response.
func (c *Client) doRequest(ctx context.Context, method, p string, query url.Values, body any, out any) error {
	req, err := c.newRequest(ctx, method, p, query, body)
	if err != nil {
		return err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the full response body so we can include it in error messages if needed.
	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// First try to decode into the generic APIResponse schema.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var apiErrResp APIResponse[json.RawMessage]
		if err := json.Unmarshal(raw, &apiErrResp); err != nil {
			// Fallback to a raw error message
			return &APIError{
				StatusCode: resp.StatusCode,
				Code:       "http_error",
				Message:    string(bytes.TrimSpace(raw)),
			}
		}
		// For better error visibility, attach the raw response body.
		msg := apiErrResp.Message
		if msg == "" {
			msg = string(bytes.TrimSpace(raw))
		} else {
			msg = fmt.Sprintf("%s; body=%s", msg, string(bytes.TrimSpace(raw)))
		}
		return &APIError{
			StatusCode: resp.StatusCode,
			Code:       apiErrResp.Error,
			Message:    msg,
		}
	}

	if out == nil {
		// Nothing to decode
		return nil
	}

	// Successful response: decode into APIResponse[T] and map Data into out.
	// out must be a pointer to the concrete type.
	wrapped := &APIResponse[json.RawMessage]{}
	if err := json.Unmarshal(raw, wrapped); err != nil {
		// Some endpoints might not use the wrapper – fallback to decoding directly into out.
		if err := json.Unmarshal(raw, out); err != nil {
			return fmt.Errorf("Antwort konnte nicht dekodiert werden: %w", err)
		}
		return nil
	}

	if !wrapped.Success && wrapped.Error != "" {
		return &APIError{
			StatusCode: resp.StatusCode,
			Code:       wrapped.Error,
			Message:    wrapped.Message,
		}
	}

	if len(wrapped.Data) == 0 {
		// No data field but success – out stays at its zero value.
		return nil
	}

	if err := json.Unmarshal(wrapped.Data, out); err != nil {
		return fmt.Errorf("Data konnte nicht dekodiert werden: %w", err)
	}

	return nil
}

