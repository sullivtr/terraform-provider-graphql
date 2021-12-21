package graphql

import (
	"context"
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
	"github.com/zclconf/go-cty/cty/function/stdlib"
	"github.com/zclconf/go-cty/cty/gocty"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"url": {
				Required:    true,
				Type:        schema.TypeString,
				ForceNew:    true,
				DefaultFunc: schema.EnvDefaultFunc("TF_GRAPHQL_URL", nil),
			},
			"headers": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
				ForceNew: true,
			},
			"oauth2_login_query": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"oauth2_login_query_variables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"oauth2_login_query_value_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"graphql_mutation": resourceGraphqlMutation(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"graphql_query": dataSourceGraphql(),
		},
		ConfigureContextFunc: graphqlConfigure,
	}
}

func graphqlConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	config := &graphqlProviderConfig{
		GQLServerUrl:   d.Get("url").(string),
		RequestHeaders: d.Get("headers").(map[string]interface{}),
	}

	if oauth2LoginQueryValueAttribute := d.Get("oauth2_login_query_value_attribute"); d.Get("oauth2_login_query") != "" && oauth2LoginQueryValueAttribute != "" {
		queryResponse, resBytes, err := queryExecute(ctx, d, config, "oauth2_login_query", "oauth2_login_query_variables")
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("unable to execute oauth2_login_query: %w", err))
		}

		if queryErrors := queryResponse.ProcessErrors(); queryErrors.HasError() {
			return nil, *queryErrors
		}

		var queryResponseCty cty.Value
		if queryResponseCty, err = gocty.ToCtyValue(string(resBytes), cty.String); err != nil {
			return nil, diag.FromErr(err)
		}

		evalCtx := &hcl.EvalContext{
			Variables: map[string]cty.Value{
				"oauth2_login_query_response": queryResponseCty,
			},
			Functions: map[string]function.Function{
				"jsondecode": stdlib.JSONDecodeFunc,
			},
		}

		var expr hcl.Expression
		var hclDiags hcl.Diagnostics
		valueTemplate := fmt.Sprintf("${jsondecode(oauth2_login_query_response).%s}", oauth2LoginQueryValueAttribute.(string))
		if expr, hclDiags = hclsyntax.ParseTemplate([]byte(valueTemplate), "", hcl.InitialPos); len(hclDiags) > 0 {
			return nil, convertDiagnosticsFromHCLToTerraformSDK(hclDiags)
		}

		var interpolatedValue cty.Value
		if interpolatedValue, hclDiags = expr.Value(evalCtx); len(hclDiags) > 0 {
			return nil, convertDiagnosticsFromHCLToTerraformSDK(hclDiags)
		}

		config.RequestAuthorizationHeaders = map[string]interface{}{
			"Authorization": fmt.Sprintf("Bearer %s", interpolatedValue.AsString()),
		}
	}

	return config, diag.Diagnostics{}
}

type graphqlProviderConfig struct {
	GQLServerUrl   string
	RequestHeaders map[string]interface{}

	RequestAuthorizationHeaders map[string]interface{}
}

func convertDiagnosticsFromHCLToTerraformSDK(hclDiags hcl.Diagnostics) diag.Diagnostics {
	diags := make(diag.Diagnostics, len(hclDiags))
	for i, hclDiag := range hclDiags {
		// HCL severity enum is: 0 (invalid), 1 (error), 2 (warning)
		// terraform SDK severity enum is: 0 (error), 1 (warning)
		diags[i].Severity = diag.Severity(hclDiag.Severity - 1)
		diags[i].Summary = hclDiag.Summary
		diags[i].Detail = hclDiag.Detail
	}
	return diags
}
