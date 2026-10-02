package gqlclient

import (
	"net/http"
)

// Client is a GraphQL client bound to a single endpoint.
type Client struct {
	Endpoint string
	Headers  map[string]interface{}
	// AuthHeaders are set after construction once an optional oauth2 login has
	// resolved. They start nil and are read-only during resource operations.
	AuthHeaders map[string]interface{}

	http *http.Client
}

type clientConfig struct {
	rateLimitPerSecond float64
	rateLimitBurst     int
}

// Option configures a Client at construction time.
type Option func(*clientConfig)

// WithRateLimit paces every request through a shared token-bucket limiter at
// perSecond requests per second with the given burst. A perSecond of 0 leaves
// rate limiting off, which is the default when the option is not supplied.
func WithRateLimit(perSecond float64, burst int) Option {
	return func(c *clientConfig) {
		c.rateLimitPerSecond = perSecond
		c.rateLimitBurst = burst
	}
}

// NewClient builds a Client for a provider instance.
func NewClient(endpoint string, headers map[string]interface{}, opts ...Option) *Client {
	cfg := clientConfig{rateLimitBurst: 1}
	for _, opt := range opts {
		opt(&cfg)
	}

	return &Client{
		Endpoint: endpoint,
		Headers:  headers,
		http:     newHTTPClient(cfg.rateLimitPerSecond, cfg.rateLimitBurst),
	}
}
