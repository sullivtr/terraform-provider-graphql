package graphql

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	return config, diag.Diagnostics{}
}

type graphqlProviderConfig struct {
	GQLServerUrl   string
	RequestHeaders map[string]interface{}
}
