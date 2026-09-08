package e2e

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const rateLimitTimestampsURL = "http://localhost:8080/rate-limit-timestamps"

// resetRateLimitWorkspace clears the shared server state file and the local
// terraform state so a subsequent apply starts from a clean slate.
func resetRateLimitWorkspace() {
	os.Remove("./gql-server/test.json")
	os.Remove("./test_rate_limit/terraform.tfstate")
	os.Remove("./test_rate_limit/terraform.tfstate.backup")
	// Drop the lock file so init regenerates it against the freshly built
	// provider; a stale checksum from a previous build breaks init.
	os.Remove("./test_rate_limit/.terraform.lock.hcl")
}

// resetRecordedTimestamps clears the server's recorded request arrival times.
func resetRecordedTimestamps(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, rateLimitTimestampsURL, nil)
	require.NoError(t, err)
	resp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	resp.Body.Close()
}

// recordedTimestamps returns the arrival times the server captured for opted-in
// requests, in arrival order.
func recordedTimestamps(t *testing.T) []time.Time {
	resp, err := http.Get(rateLimitTimestampsURL)
	require.NoError(t, err)
	defer resp.Body.Close()

	var nanos []int64
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&nanos))

	times := make([]time.Time, len(nanos))
	for i, n := range nanos {
		times[i] = time.Unix(0, n)
	}
	return times
}

// span returns the elapsed time between the first and last recorded arrival.
func span(times []time.Time) time.Duration {
	if len(times) < 2 {
		return 0
	}
	return times[len(times)-1].Sub(times[0])
}

// TestRateLimitPacesRequests proves the client-side limiter spaces the
// provider's requests. The test server records the arrival time of every
// opted-in request; the test compares an unthrottled apply against a throttled
// one.
//
//   - With rate limiting disabled (the default), the apply's requests land
//     back-to-back and the recorded arrivals span only a few milliseconds.
//   - With rate limiting enabled at ratePerSecond, the same apply's requests
//     are paced ~1/ratePerSecond apart, so the recorded arrivals span at least
//     the sum of those gaps.
func TestRateLimitPacesRequests(t *testing.T) {
	const ratePerSecond = 5.0

	resetRateLimitWorkspace()
	t.Cleanup(resetRateLimitWorkspace)
	assert.NoFileExists(t, "./gql-server/test.json")

	// Unthrottled baseline: requests are not paced, so the arrivals cluster.
	resetRecordedTimestamps(t)
	offOptions := &terraform.Options{
		TerraformDir: "./test_rate_limit",
		Vars:         map[string]interface{}{"rate_limit_per_second": 0},
		Logger:       logger.Discard,
	}
	terraform.InitAndApply(t, offOptions)
	offTimes := recordedTimestamps(t)
	terraform.Destroy(t, offOptions)
	require.GreaterOrEqual(t, len(offTimes), 3, "expected the apply to issue several opted-in requests")
	offSpan := span(offTimes)

	resetRateLimitWorkspace()

	// Throttled run: the client paces requests at ratePerSecond with burst 1.
	resetRecordedTimestamps(t)
	onOptions := &terraform.Options{
		TerraformDir: "./test_rate_limit",
		Vars: map[string]interface{}{
			"rate_limit_per_second": ratePerSecond,
			"rate_limit_burst":      1,
		},
		Logger: logger.Discard,
	}
	terraform.InitAndApply(t, onOptions)
	onTimes := recordedTimestamps(t)
	terraform.Destroy(t, onOptions)
	require.GreaterOrEqual(t, len(onTimes), 3, "expected the apply to issue several opted-in requests")
	onSpan := span(onTimes)

	// With burst 1 the first request is free and each of the remaining
	// (n-1) waits ~1/rate. Allow slack below the theoretical minimum.
	minExpectedSpan := time.Duration(float64(len(onTimes)-1)/ratePerSecond*float64(time.Second)) * 9 / 10
	assert.GreaterOrEqualf(t, onSpan, minExpectedSpan,
		"throttled apply spanned %s over %d requests, expected at least %s", onSpan, len(onTimes), minExpectedSpan)

	// The unthrottled apply must be dramatically tighter, proving the spacing
	// comes from the limiter rather than inherent server latency.
	assert.Lessf(t, offSpan, minExpectedSpan,
		"unthrottled apply spanned %s over %d requests, expected well under %s", offSpan, len(offTimes), minExpectedSpan)
}
