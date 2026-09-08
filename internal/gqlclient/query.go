package gqlclient

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ExecuteQuery runs a GraphQL query and dispatches to the paginated or single
// path. It deals only in query text and variables; extracting those from
// Terraform resource data is the caller's job.
func (c *Client) ExecuteQuery(ctx context.Context, query string, variables map[string]interface{}, usePagination bool) (*GqlQueryResponse, []byte, error) {
	if usePagination {
		return c.executePaginatedQuery(ctx, query, variables)
	}
	return c.executeSingleQuery(ctx, query, variables)
}

func prepareQueryVariables(inputVariables map[string]interface{}, cursor string) map[string]interface{} {
	currentVars := make(map[string]interface{})

	// Copy input variables
	for k, v := range inputVariables {
		js, isJS := isJSON(v)
		if isJS {
			currentVars[k] = js
		} else {
			currentVars[k] = v
		}
	}

	// Add cursor for pagination if provided
	if cursor != "" {
		currentVars["after"] = cursor
	}

	return currentVars
}

func (c *Client) executeGraphQLRequest(ctx context.Context, query string, variables map[string]interface{}) (*GqlQueryResponse, []byte, error) {
	var queryBodyBuffer bytes.Buffer

	queryObj := GqlQuery{
		Query:     query,
		Variables: variables,
	}

	if err := json.NewEncoder(&queryBodyBuffer).Encode(queryObj); err != nil {
		return nil, nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.Endpoint, &queryBodyBuffer)
	if err != nil {
		return nil, nil, err
	}

	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")
	for key, value := range c.AuthHeaders {
		req.Header.Set(key, value.(string))
	}
	for key, value := range c.Headers {
		req.Header.Set(key, value.(string))
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var gqlResponse GqlQueryResponse
	if err := json.Unmarshal(body, &gqlResponse); err != nil {
		return nil, nil, fmt.Errorf("unable to parse graphql server response: %v ---> %s", err, string(body))
	}

	return &gqlResponse, body, nil
}

func (c *Client) executeSingleQuery(ctx context.Context, query string, inputVariables map[string]interface{}) (*GqlQueryResponse, []byte, error) {
	variables := prepareQueryVariables(inputVariables, "")
	return c.executeGraphQLRequest(ctx, query, variables)
}

func (c *Client) executePaginatedQuery(ctx context.Context, query string, inputVariables map[string]interface{}) (*GqlQueryResponse, []byte, error) {
	var allResponses []GqlQueryResponse
	var finalResponseData []map[string]interface{}
	var finalResponseErrors []GqlError
	var lastCursor string

	for {
		variables := prepareQueryVariables(inputVariables, lastCursor)

		gqlResponse, _, err := c.executeGraphQLRequest(ctx, query, variables)
		if err != nil {
			return nil, nil, err
		}

		allResponses = append(allResponses, *gqlResponse)

		// Extract pageInfo from response
		pageInfo, ok := findPageInfo(gqlResponse.Data)
		if !ok {
			return nil, nil, fmt.Errorf("paginated query enabled but no pageInfo found in response (updated)")
		}

		hasNextPage, ok := pageInfo["hasNextPage"].(bool)
		if !ok {
			return nil, nil, fmt.Errorf("invalid or missing hasNextPage in pageInfo")
		}

		if !hasNextPage {
			break
		}

		endCursor, ok := pageInfo["endCursor"].(string)
		if !ok {
			return nil, nil, fmt.Errorf("invalid or missing endCursor in pageInfo")
		}
		lastCursor = endCursor
	}

	// Merge all responses
	for _, resp := range allResponses {
		// Merge the data from each response
		finalResponseData = append(finalResponseData, resp.Data)

		// Merge any errors
		finalResponseErrors = append(finalResponseErrors, resp.Errors...)
	}

	finalResponse := GqlQueryResponse{
		PaginatedResponseData: finalResponseData,
		Errors:                finalResponseErrors,
	}

	responseBytes, err := json.Marshal(finalResponse)
	if err != nil {
		return nil, nil, fmt.Errorf("error marshaling merged response: %v", err)
	}
	return &finalResponse, responseBytes, nil
}

// findPageInfo recursively searches for the "pageInfo" key in a nested map.
func findPageInfo(data map[string]interface{}) (map[string]interface{}, bool) {
	for key, value := range data {
		if key == "pageInfo" {
			if pageInfo, ok := value.(map[string]interface{}); ok {
				return pageInfo, true
			}
		}
		if nestedMap, ok := value.(map[string]interface{}); ok {
			if pageInfo, found := findPageInfo(nestedMap); found {
				return pageInfo, true
			}
		}
	}
	return nil, false
}

// isJSON checks if s can be interpreted as valid JSON, returning the
// unmarshalled value if so.
func isJSON(s interface{}) (interface{}, bool) {
	var js interface{}
	err := json.Unmarshal([]byte(s.(string)), &js)
	if err != nil {
		return nil, false
	}
	return js, true
}
