package tools

import (
	"fmt"
	"net/http"
	"net/url"
)

// GetHTTPClient returns an HTTP client configured with optional proxy support
func GetHTTPClient(proxyURL string) (*http.Client, error) {
	// If no proxy URL provided, return default client
	if proxyURL == "" {
		return http.DefaultClient, nil
	}

	// Parse proxy URL
	u, err := url.Parse(proxyURL)
	if err != nil {
		return nil, fmt.Errorf("invalid proxy URL: %w", err)
	}

	switch u.Scheme {
	case "http", "https":
		// HTTP proxy
		transport := &http.Transport{
			Proxy: http.ProxyURL(u),
		}
		return &http.Client{Transport: transport}, nil

	case "socks5":
		// SOCKS5 proxy support would require additional dependencies
		// For now, return an error with helpful message
		return nil, fmt.Errorf("SOCKS5 proxy is not currently supported. Please use HTTP proxy instead: %s", proxyURL)

	default:
		return nil, fmt.Errorf("unsupported proxy scheme: %s (supported: http, https)", u.Scheme)
	}
}
