package turvo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"drumkit-take-home/internal/config"
)

type Client struct {
	baseURL      string
	apiKey       string
	clientName   string
	clientSecret string
	username     string
	password     string
	httpClient   *http.Client

	mu          sync.Mutex
	accessToken string
	expiresAt   time.Time
}

type authRequest struct {
	GrantType string `json:"grant_type"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Scope     string `json:"scope"`
	Type      string `json:"type"`
}

type authResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

type Pagination struct {
	Start              int  `json:"start"`
	PageSize           int  `json:"pageSize"`
	TotalRecordsInPage int  `json:"totalRecordsInPage"`
	MoreAvailable      bool `json:"moreAvailable"`
}

type listShipmentsResponse struct {
	Status  string `json:"Status"`
	Details struct {
		Pagination Pagination        `json:"pagination"`
		Shipments  []json.RawMessage `json:"shipments"`
	} `json:"details"`
}

type ShipmentsResult struct {
	Count      int               `json:"count"`
	Shipments  []json.RawMessage `json:"shipments"`
	Pagination Pagination        `json:"pagination"`
}

type apiStatusError struct {
	Method string
	Path   string
	Status int
	Body   string
}

func (e *apiStatusError) Error() string {
	return fmt.Sprintf("turvo %s %s failed with status %d: %s", e.Method, e.Path, e.Status, e.Body)
}

func NewClient(cfg config.Config) *Client {
	return &Client{
		// Normalize the base URL so path concatenation doesn't produce `//...`
		// when TURVO_BASE_URL is configured with a trailing slash.
		//
		// Again, just a reasonable default.
		baseURL:      strings.TrimRight(cfg.TurvoBaseURL, "/"),
		apiKey:       cfg.TurvoAPIKey,
		clientName:   cfg.TurvoClientName,
		clientSecret: cfg.TurvoClientSecret,
		username:     cfg.TurvoUsername,
		password:     cfg.TurvoPassword,
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // Reasonable default.
		},
	}
}

func (c *Client) ListAllShipments(ctx context.Context) (ShipmentsResult, error) {
	var (
		all       []json.RawMessage
		start     *int
		lastPage  Pagination
		pageCount int
	)

	for {
		page, err := c.listShipmentsPage(ctx, start)
		if err != nil {
			return ShipmentsResult{}, err
		}

		all = append(all, page.Details.Shipments...)
		lastPage = page.Details.Pagination
		pageCount++

		if !page.Details.Pagination.MoreAvailable {
			break
		}

		next := page.Details.Pagination.Start
		start = &next

		if pageCount > 500 {
			return ShipmentsResult{}, fmt.Errorf("aborting pagination after %d pages", pageCount)
		}
	}

	return ShipmentsResult{
		Count:      len(all),
		Shipments:  all,
		Pagination: lastPage,
	}, nil
}

func (c *Client) listShipmentsPage(ctx context.Context, start *int) (listShipmentsResponse, error) {
	values := url.Values{}
	if start != nil {
		values.Set("start", strconv.Itoa(*start))
	}

	path := "/shipments/list"
	if encoded := values.Encode(); encoded != "" {
		path += "?" + encoded
	}

	var response listShipmentsResponse
	if err := c.getJSON(ctx, path, &response); err != nil {
		return listShipmentsResponse{}, err
	}

	return response, nil
}

func (c *Client) getJSON(ctx context.Context, path string, out any) error {
	return c.sendJSON(ctx, http.MethodGet, path, nil, out)
}

func (c *Client) postJSON(ctx context.Context, path string, requestBody any, out any) error {
	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("marshal request body: %w", err)
	}

	return c.sendJSON(ctx, http.MethodPost, path, body, out)
}

func (c *Client) sendJSON(ctx context.Context, method, path string, requestBody []byte, out any) error {
	token, err := c.token(ctx)
	if err != nil {
		return err
	}

	status, body, err := c.do(ctx, method, path, token, requestBody)
	if err != nil {
		return err
	}

	if status == http.StatusUnauthorized {
		c.invalidateToken()
		token, err = c.token(ctx)
		if err != nil {
			return err
		}

		status, body, err = c.do(ctx, method, path, token, requestBody)
		if err != nil {
			return err
		}
	}

	if status < 200 || status >= 300 {
		return &apiStatusError{
			Method: method,
			Path:   path,
			Status: status,
			Body:   strings.TrimSpace(string(body)),
		}
	}

	if out == nil {
		return nil
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("decode turvo response: %w", err)
	}

	return nil
}

func (c *Client) token(ctx context.Context) (string, error) {
	c.mu.Lock()
	if c.accessToken != "" && time.Now().Before(c.expiresAt.Add(-1*time.Minute)) {
		token := c.accessToken
		c.mu.Unlock()
		return token, nil
	}
	c.mu.Unlock()

	return c.authenticate(ctx)
}

func (c *Client) invalidateToken() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.accessToken = ""
	c.expiresAt = time.Time{}
}

func (c *Client) authenticate(ctx context.Context) (string, error) {
	c.mu.Lock()
	if c.accessToken != "" && time.Now().Before(c.expiresAt.Add(-1*time.Minute)) {
		token := c.accessToken
		c.mu.Unlock()
		return token, nil
	}
	c.mu.Unlock()

	payload := authRequest{
		GrantType: "password",
		Username:  c.username,
		Password:  c.password,
		Scope:     "read+trust+write",
		Type:      "business",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("marshal auth request: %w", err)
	}

	query := url.Values{}
	query.Set("client_id", c.clientName)
	query.Set("client_secret", c.clientSecret)

	status, responseBody, err := c.do(ctx, http.MethodPost, "/oauth/token?"+query.Encode(), "", body)
	if err != nil {
		return "", err
	}
	if status < 200 || status >= 300 {
		return "", fmt.Errorf("turvo auth failed with status %d: %s", status, strings.TrimSpace(string(responseBody)))
	}

	var response authResponse
	if err := json.Unmarshal(responseBody, &response); err != nil {
		return "", fmt.Errorf("decode auth response: %w", err)
	}
	if response.AccessToken == "" {
		return "", fmt.Errorf("turvo auth response missing access_token")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.accessToken = response.AccessToken
	if response.ExpiresIn > 0 {
		c.expiresAt = time.Now().Add(time.Duration(response.ExpiresIn) * time.Second)
	} else {
		c.expiresAt = time.Now().Add(12 * time.Hour)
	}

	return c.accessToken, nil
}

func (c *Client) do(ctx context.Context, method, path, bearerToken string, body []byte) (int, []byte, error) {
	requestURL := c.baseURL + path

	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, requestURL, reader)
	if err != nil {
		return 0, nil, fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("x-api-key", c.apiKey)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if bearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+bearerToken)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return 0, nil, fmt.Errorf("perform request %s %s: %w", method, requestURL, err)
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, nil, fmt.Errorf("read response body: %w", err)
	}

	return resp.StatusCode, responseBody, nil
}
