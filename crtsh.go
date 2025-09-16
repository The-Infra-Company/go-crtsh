// Package gocrtsh provides a client for the crt.sh API
package gocrtsh

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	baseURL   = "https://crt.sh"
	version   = "0.1.0"
	userAgent = "go-crtsh/" + version
)

// CertificateRecord represents a record in the crt.sh database
type CertificateRecord struct {
	IssuerCAID     int64  `json:"issuer_ca_id"`
	IssuerName     string `json:"issuer_name"`
	CommonName     string `json:"common_name"`
	NameValue      string `json:"name_value"`
	ID             int64  `json:"id"`
	EntryTimestamp string `json:"entry_timestamp"`
	NotBefore      string `json:"not_before"`
	NotAfter       string `json:"not_after"`
	SerialNumber   string `json:"serial_number"`
	ResultCount    int64  `json:"result_count"`
}

// Client is the main handler for crt.sh API interactions
type Client struct {
	client  *http.Client
	baseURL string
}

// New creates a new crt.sh API client with default settings
func New() *Client {
	return &Client{
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
	}
}

// NewWithClient creates a new crt.sh API client with a custom HTTP client
func NewWithClient(httpClient *http.Client) *Client {
	return &Client{
		client:  httpClient,
		baseURL: baseURL,
	}
}

// SearchOptions contains options for certificate searches
type SearchOptions struct {
	// IncludeExpired determines whether to include expired certificates
	IncludeExpired bool
	// Wildcard determines whether to search for subdomains (adds %. prefix if not present)
	Wildcard bool
}

// Search searches for certificates for the given domain
func (c *Client) Search(ctx context.Context, domain string, opts *SearchOptions) ([]CertificateRecord, error) {
	if domain == "" {
		return nil, fmt.Errorf("domain cannot be empty")
	}

	// Validate basic domain format
	if strings.Contains(domain, " ") || strings.HasPrefix(domain, ".") || strings.HasSuffix(domain, ".") {
		return nil, fmt.Errorf("invalid domain format: %s", domain)
	}

	searchDomain := domain
	if opts != nil && opts.Wildcard && !strings.Contains(domain, "%") && !strings.HasPrefix(domain, "*.") {
		searchDomain = "%." + domain
	}

	apiURL, err := c.buildURL(searchDomain, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to build URL: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "GET", apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Accept", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Handle empty response
	if len(body) == 0 || string(body) == "[]" {
		return []CertificateRecord{}, nil
	}

	var records []CertificateRecord
	if err := json.Unmarshal(body, &records); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return records, nil
}

// buildURL constructs the API URL with proper query parameters
func (c *Client) buildURL(domain string, opts *SearchOptions) (string, error) {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return "", err
	}

	// Add the path if not present
	if u.Path == "" {
		u.Path = "/"
	}

	// Build query parameters
	params := url.Values{}
	params.Set("q", domain)
	params.Set("output", "json")

	// Add exclude expired if specified
	if opts != nil && !opts.IncludeExpired {
		params.Set("exclude", "expired")
	}

	u.RawQuery = params.Encode()
	return u.String(), nil
}

// BasicSearch is a convenience method for basic searches
func (c *Client) BasicSearch(ctx context.Context, domain string) ([]CertificateRecord, error) {
	return c.Search(ctx, domain, nil)
}

// SearchWithWildcard is a convenience method for wildcard searches
func (c *Client) SearchWithWildcard(ctx context.Context, domain string, includeExpired bool) ([]CertificateRecord, error) {
	opts := &SearchOptions{
		IncludeExpired: includeExpired,
		Wildcard:       true,
	}
	return c.Search(ctx, domain, opts)
}
