package gqlclient

import (
	"bytes"
	"log"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

// captureLog redirects the std logger for the duration of fn and returns what
// was written, restoring the previous output afterwards.
func captureLog(t *testing.T, fn func()) string {
	t.Helper()
	var buf bytes.Buffer
	prevOut, prevFlags := log.Writer(), log.Flags()
	log.SetOutput(&buf)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(prevOut)
		log.SetFlags(prevFlags)
	})
	fn()
	return buf.String()
}

const tracingLog = "Enabling HTTP requests/responses tracing"

func TestNewHTTPClientTransport(t *testing.T) {
	// IsDebugOrHigher is true for DEBUG and TRACE, false otherwise; the debug
	// path wraps the transport and logs a breadcrumb, the quiet path does not.
	tests := []struct {
		tfLog     string
		wantDebug bool
	}{
		{"TRACE", true},
		{"DEBUG", true},
		{"INFO", false},
		{"WARN", false},
		{"ERROR", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.tfLog, func(t *testing.T) {
			t.Setenv("TF_LOG", tt.tfLog)

			var client *http.Client
			out := captureLog(t, func() { client = newHTTPClient() })

			logged := strings.Contains(out, tracingLog)
			if logged != tt.wantDebug {
				t.Fatalf("tracing log present = %v, want %v (out %q)", logged, tt.wantDebug, out)
			}

			wrapped := client.Transport != http.DefaultTransport
			if wrapped != tt.wantDebug {
				t.Fatalf("transport wrapped = %v, want %v (got %T)", wrapped, tt.wantDebug, client.Transport)
			}

			if tt.wantDebug {
				// logging.NewTransport returns an unexported *logging.transport;
				// assert by type name so a future rename surfaces here.
				if name := reflect.TypeOf(client.Transport).String(); name != "*logging.transport" {
					t.Fatalf("expected *logging.transport wrapper, got %s", name)
				}
			}
		})
	}
}
