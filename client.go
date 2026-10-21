package platform5

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const defaultBaseURL = "https://api.advantaplus.ai"

// Client is the Platform5 Developer API client.
// Create one with [New], then call methods on the service accessors
// (e.g. [Client.SMS], [Client.Email], [Client.Messages], [Client.Account]).
type Client struct {
	apiKey  string
	baseURL string
	http    *http.Client

	sms      *smsService
	email    *emailService
	messages *messagesService
	account  *accountService
}

// New creates a new Platform5 API client.
//
// The client lazily initialises each service on first access, so
// creating a client is cheap regardless of how many services are used.
func New(cfg Config) *Client {
	baseURL := cfg.BaseURL
	if baseURL == "" {
		baseURL = defaultBaseURL
	}
	return &Client{
		apiKey:  cfg.APIKey,
		baseURL: strings.TrimRight(baseURL, "/"),
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// SMS returns the SMS service.
func (c *Client) SMS() *smsService {
	if c.sms == nil {
		c.sms = &smsService{client: c}
	}
	return c.sms
}

// Email returns the Email service.
func (c *Client) Email() *emailService {
	if c.email == nil {
		c.email = &emailService{client: c}
	}
	return c.email
}

// Messages returns the Messages service for checking delivery status.
func (c *Client) Messages() *messagesService {
	if c.messages == nil {
		c.messages = &messagesService{client: c}
	}
	return c.messages
}

// Account returns the Account service for checking balance.
func (c *Client) Account() *accountService {
	if c.account == nil {
		c.account = &accountService{client: c}
	}
	return c.account
}

// Health checks whether the API is reachable.
// It returns nil on success or an error describing the failure.
func (c *Client) Health(ctx context.Context) error {
	_, err := c.request(ctx, http.MethodGet, "/health", nil, "")
	return err
}

func (c *Client) request(ctx context.Context, method, path string, body any, idempotencyKey string) (json.RawMessage, error) {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-Key", c.apiKey)
	if idempotencyKey != "" {
		req.Header.Set("Idempotency-Key", idempotencyKey)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	requestID := resp.Header.Get("X-Request-ID")

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		var env apiResponse[json.RawMessage]
		if err := json.Unmarshal(respBody, &env); err != nil {
			env.Message = http.StatusText(resp.StatusCode)
		}
		return nil, c.toError(resp.StatusCode, resp.Header, &env, requestID)
	}

	return respBody, nil
}

func (c *Client) toError(status int, h http.Header, body *apiResponse[json.RawMessage], requestID string) error {
	err := &Error{
		StatusCode: status,
		Message:    body.Message,
		Errors:     body.Errors,
		RequestID:  requestID,
	}

	if status == 429 {
		limit, _ := strconv.Atoi(h.Get("X-RateLimit-Limit"))
		remaining, _ := strconv.Atoi(h.Get("X-RateLimit-Remaining"))
		return &RateLimitError{
			Err:       *err,
			Limit:     limit,
			Remaining: remaining,
		}
	}

	return err
}

// decode unmarshals the standard API response envelope and returns the data field.
func decode[T any](raw json.RawMessage) (T, error) {
	var env apiResponse[T]
	if err := json.Unmarshal(raw, &env); err != nil {
		var zero T
		return zero, fmt.Errorf("decode response: %w", err)
	}
	return env.Data, nil
}
