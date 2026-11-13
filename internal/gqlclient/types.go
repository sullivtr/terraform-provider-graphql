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
