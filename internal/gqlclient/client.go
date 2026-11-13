package gqlclient

import (
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
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

// NewClient builds a Client for a provider instance.
func NewClient(endpoint string, headers map[string]interface{}) *Client {
	return &Client{
		Endpoint: endpoint,
		Headers:  headers,
		http:     newHTTPClient(),
	}
}

func newHTTPClient() *http.Client {
	transport := http.DefaultTransport
	if logging.IsDebugOrHigher() {
		log.Printf("[DEBUG] Enabling HTTP requests/responses tracing")
		transport = logging.NewTransport("GraphQL", transport)
	}
	return &http.Client{Transport: transport}
}
