package gqlclient

import (
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
)

// Client is a GraphQL client bound to a single endpoint. It owns the static
// and authorization headers and a shared *http.Client. ExecuteQuery is its only
// public behavior.
type Client struct {
	// Endpoint is the GraphQL server URL. GraphQL is single-endpoint: the query
	// body — not the path — drives every request, so this never changes per call.
	Endpoint string
	// Headers are static request headers supplied by the provider config.
	Headers map[string]interface{}
	// AuthHeaders are set after construction once an optional oauth2 login has
	// resolved. They start nil and are read-only during resource operations.
	AuthHeaders map[string]interface{}

	http *http.Client
}

// NewClient builds a Client for a provider instance. AuthHeaders start empty;
// the caller sets them after any oauth2 login query resolves.
func NewClient(endpoint string, headers map[string]interface{}) *Client {
	return &Client{
		Endpoint: endpoint,
		Headers:  headers,
		http:     newHTTPClient(),
	}
}

// newHTTPClient builds the shared http.Client once per provider instance,
// wiring the debug tracing transport when TF logging is at DEBUG or higher.
func newHTTPClient() *http.Client {
	transport := http.DefaultTransport
	if logging.IsDebugOrHigher() {
		log.Printf("[DEBUG] Enabling HTTP requests/responses tracing")
		transport = logging.NewTransport("GraphQL", transport)
	}
	return &http.Client{Transport: transport}
}
