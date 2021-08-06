package graphql

import "github.com/hashicorp/terraform-plugin-sdk/v2/diag"

type GqlQuery struct {
	Query     string                 `json:"query,omitempty"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GqlQueryResponse struct {
	Data   map[string]interface{} `json:"data,omitempty"`
	Errors []GqlError             `json:"errors,omitempty"`
}

type GqlError struct {
	Message string `json:"message,omitempty"`
}

func (r *GqlQueryResponse) ProcessErrors() *diag.Diagnostics {
	var diags diag.Diagnostics
	if r.Errors != nil && len(r.Errors) > 0 {
		for _, queryErr := range r.Errors {
			diags = append(diags, diag.Diagnostic{Severity: diag.Error, Detail: queryErr.Message})
		}
	}
	return &diags
}
