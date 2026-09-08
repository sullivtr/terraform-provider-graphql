package gqlclient

import (
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"golang.org/x/time/rate"
)

// rateLimitedTransport wraps a base http.RoundTripper and blocks on a shared
// token-bucket limiter before delegating each request. Wrapping rather than
// replacing lets it stack with the logging.NewTransport debug wrapper.
type rateLimitedTransport struct {
	base    http.RoundTripper
	limiter *rate.Limiter
}

func (t *rateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Wait blocks until a token is free or the request context is cancelled.
	if err := t.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}
	return t.base.RoundTrip(req)
}

// newHTTPClient builds the shared http.Client for a provider instance, wiring
// the debug tracing transport when TF logging is at DEBUG or higher. When
// ratePerSecond is 0 the limiter wrapper is skipped entirely, so the default
// path is unchanged.
func newHTTPClient(ratePerSecond float64, burst int) *http.Client {
	base := http.DefaultTransport
	if logging.IsDebugOrHigher() {
		log.Printf("[DEBUG] Enabling HTTP requests/responses tracing")
		base = logging.NewTransport("GraphQL", base)
	}

	if ratePerSecond > 0 {
		base = &rateLimitedTransport{
			base:    base,
			limiter: rate.NewLimiter(rate.Limit(ratePerSecond), burst),
		}
	}

	return &http.Client{Transport: base}
}
