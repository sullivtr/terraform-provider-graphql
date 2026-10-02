package graphql

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/sullivtr/terraform-provider-graphql/internal/gqlclient"
)

// processErrors turns GraphQL server errors into Terraform diagnostics. It lives
// in the provider package rather than on the response type so gqlclient stays
// free of the Terraform SDK's diag types.
func processErrors(r *gqlclient.GqlQueryResponse) *diag.Diagnostics {
	var diags diag.Diagnostics
	for _, queryErr := range r.Errors {
		msg := fmt.Sprintf("graphql server error: %s", queryErr.Message)
		diags = append(diags, diag.Diagnostic{Summary: msg, Severity: diag.Error, Detail: msg})
	}
	return &diags
}
