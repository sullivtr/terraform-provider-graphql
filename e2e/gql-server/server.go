package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/sullivtr/terraform-provider-graphql/gql-server/graph"
	"github.com/sullivtr/terraform-provider-graphql/gql-server/graph/generated"
)

const defaultPort = "8080"

// rateLimitHeader opts a request into the server's arrival-time recording. Only
// the rate limiting e2e fixture sends it, so the other fixtures are unaffected.
const rateLimitHeader = "x-e2e-rate-limit"

// requestRecorder captures the arrival time of every opted-in /query request so
// a test can assert how far apart the client's paced requests landed.
type requestRecorder struct {
	mu    sync.Mutex
	times []time.Time
}

func (r *requestRecorder) record() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.times = append(r.times, time.Now())
}

func (r *requestRecorder) reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.times = nil
}

// snapshotNanos returns the recorded arrival times as Unix-nanosecond values,
// in the order they arrived, for a test to inspect the spacing.
func (r *requestRecorder) snapshotNanos() []int64 {
	r.mu.Lock()
	defer r.mu.Unlock()
	out := make([]int64, len(r.times))
	for i, t := range r.times {
		out[i] = t.UnixNano()
	}
	return out
}

// recordMiddleware timestamps opted-in requests before passing them to the
// GraphQL handler. Non-opted-in traffic is untouched.
func recordMiddleware(rec *requestRecorder, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get(rateLimitHeader) != "" {
			rec.record()
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	rec := &requestRecorder{}

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", recordMiddleware(rec, srv))

	// Test-only introspection of the recorded request arrival times. GET returns
	// the timestamps as JSON; DELETE clears them so each test starts clean.
	http.HandleFunc("/rate-limit-timestamps", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			rec.reset()
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(rec.snapshotNanos())
	})

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
