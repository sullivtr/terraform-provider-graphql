package gqlclient

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExecuteQuerySendsRequest(t *testing.T) {
	var gotMethod, gotContentType, gotAccept, gotCustom, gotAuth string
	var gotBody GqlQuery

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotContentType = r.Header.Get("Content-Type")
		gotAccept = r.Header.Get("Accept")
		gotCustom = r.Header.Get("X-Custom")
		gotAuth = r.Header.Get("Authorization")

		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)

		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"data":{"greeting":"response-body-value"}}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, map[string]interface{}{"X-Custom": "custom-value"})
	client.AuthHeaders = map[string]interface{}{"Authorization": "Bearer token123"}

	resp, _, err := client.ExecuteQuery(context.Background(), "query { greeting }", map[string]interface{}{"name": "request-var-value"}, false)
	if err != nil {
		t.Fatalf("ExecuteQuery returned error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotContentType != "application/json; charset=utf-8" {
		t.Errorf("Content-Type = %q", gotContentType)
	}
	if gotAccept != "application/json; charset=utf-8" {
		t.Errorf("Accept = %q", gotAccept)
	}
	if gotCustom != "custom-value" {
		t.Errorf("X-Custom = %q, want custom-value", gotCustom)
	}
	if gotAuth != "Bearer token123" {
		t.Errorf("Authorization = %q, want Bearer token123", gotAuth)
	}
	if gotBody.Query != "query { greeting }" {
		t.Errorf("body query = %q", gotBody.Query)
	}
	if gotBody.Variables["name"] != "request-var-value" {
		t.Errorf("body variables[name] = %v, want request-var-value", gotBody.Variables["name"])
	}
	if resp.Data["greeting"] != "response-body-value" {
		t.Errorf("resp.Data[greeting] = %v, want response-body-value", resp.Data["greeting"])
	}
}

func TestExecuteQueryParsesServerErrors(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"errors":[{"message":"boom"}]}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, nil)

	resp, _, err := client.ExecuteQuery(context.Background(), "query { x }", map[string]interface{}{}, false)
	if err != nil {
		t.Fatalf("transport error: %v", err)
	}
	if len(resp.Errors) != 1 || resp.Errors[0].Message != "boom" {
		t.Fatalf("expected one server error 'boom', got %+v", resp.Errors)
	}
}

func TestExecuteQueryInvalidJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`not json`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, nil)

	if _, _, err := client.ExecuteQuery(context.Background(), "query { x }", map[string]interface{}{}, false); err == nil {
		t.Fatal("expected error on non-JSON response, got nil")
	}
}

func TestExecuteQueryPaginationFollowsCursor(t *testing.T) {
	var cursors []interface{}
	call := 0

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body GqlQuery
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &body)
		cursors = append(cursors, body.Variables["after"])

		// First page has a next page, second page ends pagination.
		if call == 0 {
			call++
			_, _ = w.Write([]byte(`{"data":{"items":{"pageInfo":{"hasNextPage":true,"endCursor":"cursor-1"}}}}`))
			return
		}
		_, _ = w.Write([]byte(`{"data":{"items":{"pageInfo":{"hasNextPage":false,"endCursor":""}}}}`))
	}))
	defer srv.Close()

	client := NewClient(srv.URL, nil)

	resp, _, err := client.ExecuteQuery(context.Background(), "query { items }", map[string]interface{}{}, true)
	if err != nil {
		t.Fatalf("ExecuteQuery paginated error: %v", err)
	}

	if len(cursors) != 2 {
		t.Fatalf("expected 2 requests, got %d", len(cursors))
	}
	if cursors[0] != nil {
		t.Errorf("first request should carry no cursor, got %v", cursors[0])
	}
	if cursors[1] != "cursor-1" {
		t.Errorf("second request cursor = %v, want cursor-1", cursors[1])
	}
	if len(resp.PaginatedResponseData) != 2 {
		t.Errorf("merged pages = %d, want 2", len(resp.PaginatedResponseData))
	}
}

func TestExecuteQueryPaginationInvalidPageInfo(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"no pageInfo", `{"data":{"items":{"nodes":[]}}}`},
		{"hasNextPage wrong type", `{"data":{"items":{"pageInfo":{"hasNextPage":"nope"}}}}`},
		{"endCursor wrong type", `{"data":{"items":{"pageInfo":{"hasNextPage":true,"endCursor":123}}}}`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				_, _ = w.Write([]byte(tt.body))
			}))
			defer srv.Close()

			client := NewClient(srv.URL, nil)

			if _, _, err := client.ExecuteQuery(context.Background(), "query { items }", map[string]interface{}{}, true); err == nil {
				t.Fatalf("expected error for %s, got nil", tt.name)
			}
		})
	}
}

func TestExecuteQueryTransportError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	url := srv.URL
	srv.Close() // server down: http.Do fails

	client := NewClient(url, nil)

	if _, _, err := client.ExecuteQuery(context.Background(), "query { x }", map[string]interface{}{}, false); err == nil {
		t.Fatal("expected transport error against closed server, got nil")
	}
}

// Variable values arrive as strings (Terraform TypeString maps). A value that
// parses as JSON is decoded into structured data before it is sent; anything
// that does not parse is forwarded as the original string.
func TestExecuteQueryEncodesVariables(t *testing.T) {
	tests := []struct {
		name  string
		input string
		check func(t *testing.T, got interface{})
	}{
		{
			name:  "plain string stays string",
			input: "just-a-string",
			check: func(t *testing.T, got interface{}) {
				if got != "just-a-string" {
					t.Errorf("got %v (%T), want unchanged string", got, got)
				}
			},
		},
		{
			name:  "json string decodes to structured data",
			input: `{"tags":["a","b"]}`,
			check: func(t *testing.T, got interface{}) {
				// We only care that it was decoded, not forwarded as the raw
				// string; the exact shape is encoding/json's contract, not ours.
				if _, ok := got.(map[string]interface{}); !ok {
					t.Errorf("got %v (%T), want decoded structured data", got, got)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotBody GqlQuery
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				raw, _ := io.ReadAll(r.Body)
				_ = json.Unmarshal(raw, &gotBody)
				_, _ = w.Write([]byte(`{"data":{"ok":true}}`))
			}))
			defer srv.Close()

			client := NewClient(srv.URL, nil)

			vars := map[string]interface{}{"v": tt.input}
			if _, _, err := client.ExecuteQuery(context.Background(), "query { x }", vars, false); err != nil {
				t.Fatalf("ExecuteQuery returned error: %v", err)
			}

			tt.check(t, gotBody.Variables["v"])
		})
	}
}
