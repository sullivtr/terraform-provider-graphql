// Package gqlclient is a GraphQL client that is independent of Terraform's
// resource/schema layer: it deals in query text and variables, never in
// *schema.ResourceData. (It does use the SDK's logging helper for debug
// tracing.) ExecuteQuery is the only entry point callers need — the
// single-request, pagination, and transport helpers are unexported so the
// package boundary enforces that.
package gqlclient

// GqlQuery is the request body posted to the GraphQL endpoint.
type GqlQuery struct {
	Query     string                 `json:"query,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GqlQueryResponse is the decoded GraphQL response. PaginatedResponseData holds
// the merged pages when a paginated query is run.
type GqlQueryResponse struct {
	Data                  map[string]interface{}   `json:"data,omitempty"`
	Errors                []GqlError               `json:"errors,omitempty"`
	PaginatedResponseData []map[string]interface{} `json:"paginatedResponseData,omitempty"`
}

// GqlError is a single error entry returned by the GraphQL server.
type GqlError struct {
	Message string `json:"message,omitempty"`
}
