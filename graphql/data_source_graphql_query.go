package graphql

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceGraphql() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"query": {
				Type:     schema.TypeString,
				Required: true,
			},
			"query_variables": {
				Type: schema.TypeMap,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"query_response": {
				Type:        schema.TypeString,
				Description: "The raw body of the HTTP response from the last read of the object.",
				Computed:    true,
			},
		},
		ReadContext: dataSourceGraphqlQuery,
	}
}

func dataSourceGraphqlQuery(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	queryResponse, resBytes, err := queryExecute(ctx, d, m, "query", "query_variables")
	if err != nil {
		return diag.FromErr(err)
	}

	if queryErrors := queryResponse.ProcessErrors(); queryErrors.HasError() {
		return *queryErrors
	}

	objID := hash(resBytes)
	d.SetId(fmt.Sprint(objID))
	if err := d.Set("query_response", string(resBytes)); err != nil {
		return diag.FromErr(err)
	}
	return diags
}

// fake
