package gqlclient

import (
	"log"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/logging"
	"golang.org/x/time/rate"
)

// rateLimitedTransport wraps rather than replaces the base RoundTripper so it
// stacks with the logging.NewTransport debug wrapper.
type rateLimitedTransport struct {
	base    http.RoundTripper
	limiter *rate.Limiter
}

func (t *rateLimitedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if err := t.limiter.Wait(req.Context()); err != nil {
		return nil, err
	}
	return t.base.RoundTrip(req)
}

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
